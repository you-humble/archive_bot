package audios

import (
	"context"
	"sync"

	"archive_bot/pkg/er"
	"archive_bot/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoAudio = er.New("there's no saved audios", "", nil)

var (
	instance *pgRepository
	once     sync.Once
)

type pgRepository struct {
	log *logger.Logger
	db  *pgxpool.Pool
}

// NewRepository creates new audios repository.
func NewRepository(ctx context.Context, log *logger.Logger, db *pgxpool.Pool) (*pgRepository, error) {
	once.Do(func() {
		instance = &pgRepository{log: log, db: db}
	})

	return instance, nil
}

// Save saves note to database.
func (repo *pgRepository) Save(ctx context.Context, n *Audio) (int, error) {
	const op string = "audios.repository.Save"
	log := repo.log.With(logger.String("operation", op))

	var id int
	if err := repo.db.QueryRow(ctx,
		`INSERT INTO audios (texts_id, file_id, media_group_id)
		 VALUES ($1, $2, $3) RETURNING id;`,
		n.TextsID, n.FileID, n.MediaGroupID).Scan(&id); err != nil {
		log.Error("save a Audio", logger.ErrAttr(err))
		return 0, er.New("unable to save a audio", op, err)
	}

	return id, nil
}

func (repo *pgRepository) FindByTextsID(ctx context.Context, id int) ([]string, error) {
	const op string = "audios.repository.FindByTextsID"

	var audio Audio
	if err := repo.db.QueryRow(ctx,
		`SELECT file_id, media_group_id FROM audios WHERE texts_id = $1;`,
		id).Scan(&audio.FileID, &audio.MediaGroupID); err != nil {
		return []string{}, er.New("the audio could not be found", op, err)
	}

	if audio.MediaGroupID != "" {
		rows, err := repo.db.Query(ctx,
			`SELECT file_id FROM audios WHERE media_group_id = $1;`,
			audio.MediaGroupID)
		if err != nil {
			return []string{}, er.New("the Audio could not be found", op, err)
		}
		defer rows.Close()

		var audioIDs []string
		for rows.Next() {
			var AudioID string
			if err := rows.Scan(&AudioID); err != nil {
				return nil, er.New("unable to scan data", op, err)
			}
			audioIDs = append(audioIDs, AudioID)
		}

		if err := rows.Err(); err != nil {
			return nil, er.New("error in rows: %w", op, err)
		}

		return audioIDs, nil
	}

	return []string{audio.FileID}, nil
}

func (repo *pgRepository) UpdateByTextsID(ctx context.Context, n *Audio) error {
	const op string = "audios.repository.UpdateByID"

	if _, err := repo.db.Exec(ctx,
		`UPDATE audios
		SET file_id = $1, media_group_id = $2
		WHERE texts_id = $3;`,
		n.FileID, n.MediaGroupID, n.TextsID); err != nil {
		return er.New("unable to move audio", op, err)
	}

	return nil
}
