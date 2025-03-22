package videos

import (
	"context"
	"sync"

	"archive_bot/pkg/er"
	"archive_bot/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoVideo = er.New("there's no saved videos", "", nil)

var (
	instance *pgRepository
	once     sync.Once
)

type pgRepository struct {
	log *logger.Logger
	db  *pgxpool.Pool
}

// NewRepository creates new videos repository.
func NewRepository(ctx context.Context, log *logger.Logger, db *pgxpool.Pool) (*pgRepository, error) {
	once.Do(func() {
		instance = &pgRepository{log: log, db: db}
	})

	return instance, nil
}

// Save saves note to database.
func (repo *pgRepository) Save(ctx context.Context, n *Video) (int, error) {
	const op string = "videos.repository.Save"
	log := repo.log.With(logger.String("operation", op))

	var id int
	if err := repo.db.QueryRow(ctx,
		`INSERT INTO videos (texts_id, file_id, media_group_id)
		 VALUES ($1, $2, $3) RETURNING id;`,
		n.TextsID, n.FileID, n.MediaGroupID).Scan(&id); err != nil {
		log.Error("save a video", logger.ErrAttr(err))
		return 0, er.New("unable to save a video", op, err)
	}

	return id, nil
}

func (repo *pgRepository) FindByTextsID(ctx context.Context, id int) ([]string, error) {
	const op string = "videos.repository.FindByTextsID"

	var video Video
	if err := repo.db.QueryRow(ctx,
		`SELECT file_id, media_group_id FROM videos WHERE texts_id = $1;`,
		id).Scan(&video.FileID, &video.MediaGroupID); err != nil {
		return []string{}, er.New("the video could not be found", op, err)
	}

	if video.MediaGroupID != "" {
		rows, err := repo.db.Query(ctx,
			`SELECT file_id FROM videos WHERE media_group_id = $1;`,
			video.MediaGroupID)
		if err != nil {
			return []string{}, er.New("the video could not be found", op, err)
		}
		defer rows.Close()

		var videoIDs []string
		for rows.Next() {
			var videoID string
			if err := rows.Scan(&videoID); err != nil {
				return nil, er.New("unable to scan data", op, err)
			}
			videoIDs = append(videoIDs, videoID)
		}

		if err := rows.Err(); err != nil {
			return nil, er.New("error in rows: %w", op, err)
		}

		return videoIDs, nil
	}

	return []string{video.FileID}, nil
}

func (repo *pgRepository) UpdateByTextsID(ctx context.Context, n *Video) error {
	const op string = "videos.repository.UpdateByID"

	if _, err := repo.db.Exec(ctx,
		`UPDATE videos
		SET file_id = $1, media_group_id = $2
		WHERE texts_id = $3;`,
		n.FileID, n.MediaGroupID, n.TextsID); err != nil {
		return er.New("unable to move video", op, err)
	}

	return nil
}
