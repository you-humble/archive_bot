package photos

import (
	"context"
	"sync"

	"archive_bot/pkg/er"
	"archive_bot/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoPhoto = er.New("there's no saved photos", "", nil)

var (
	instance *pgRepository
	once     sync.Once
)

type pgRepository struct {
	log *logger.Logger
	db  *pgxpool.Pool
}

// NewRepository creates new photos repository.
func NewRepository(ctx context.Context, log *logger.Logger, db *pgxpool.Pool) (*pgRepository, error) {
	once.Do(func() {
		instance = &pgRepository{log: log, db: db}
	})

	return instance, nil
}

// Save saves note to database.
func (repo *pgRepository) Save(ctx context.Context, n *Photo) (int, error) {
	const op string = "photos.repository.Save"
	log := repo.log.With(logger.String("operation", op))

	var id int
	if err := repo.db.QueryRow(ctx,
		`INSERT INTO photos (texts_id, file_id, media_group_id)
		 VALUES ($1, $2, $3) RETURNING id;`,
		n.TextsID, n.FileID, n.MediaGroupID).Scan(&id); err != nil {
		log.Error("save a photo", logger.ErrAttr(err))
		return 0, er.New("unable to save a photo", op, err)
	}

	return id, nil
}

func (repo *pgRepository) FindByTextsID(ctx context.Context, id int) ([]string, error) {
	const op string = "photos.repository.FindByTextsID"

	var photo Photo
	if err := repo.db.QueryRow(ctx,
		`SELECT file_id, media_group_id FROM photos WHERE texts_id = $1;`,
		id).Scan(&photo.FileID, &photo.MediaGroupID); err != nil {
		return []string{}, er.New("the photo could not be found", op, err)
	}

	if photo.MediaGroupID != "" {
		rows, err := repo.db.Query(ctx,
			`SELECT file_id FROM photos WHERE media_group_id = $1;`,
			photo.MediaGroupID)
		if err != nil {
			return []string{}, er.New("the photo could not be found", op, err)
		}
		defer rows.Close()

		var photoIDs []string
		for rows.Next() {
			var photoID string
			if err := rows.Scan(&photoID); err != nil {
				return nil, er.New("unable to scan data", op, err)
			}
			photoIDs = append(photoIDs, photoID)
		}

		if err := rows.Err(); err != nil {
			return nil, er.New("error in rows: %w", op, err)
		}

		return photoIDs, nil
	}

	return []string{photo.FileID}, nil
}

func (repo *pgRepository) UpdateByTextsID(ctx context.Context, n *Photo) error {
	const op string = "photos.repository.UpdateByID"

	if _, err := repo.db.Exec(ctx,
		`UPDATE photos
		SET file_id = $1, media_group_id = $2
		WHERE texts_id = $3;`,
		n.FileID, n.MediaGroupID, n.TextsID); err != nil {
		return er.New("unable to move photo", op, err)
	}

	return nil
}
