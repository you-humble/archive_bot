package animations

import (
	"context"
	"sync"

	"archive_bot/pkg/er"
	"archive_bot/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoAnimations = er.New("there's no saved animations", "", nil)

var (
	instance *pgRepository
	once     sync.Once
)

type pgRepository struct {
	log *logger.Logger
	db  *pgxpool.Pool
}

// NewRepository creates new animations repository.
func NewRepository(ctx context.Context, log *logger.Logger, db *pgxpool.Pool) (*pgRepository, error) {
	once.Do(func() {
		instance = &pgRepository{log: log, db: db}
	})

	return instance, nil
}

// Save saves note to database.
func (repo *pgRepository) Save(ctx context.Context, n *Animations) (int, error) {
	const op string = "animations.repository.Save"
	log := repo.log.With(logger.String("operation", op))

	var id int
	if err := repo.db.QueryRow(ctx,
		`INSERT INTO animations (texts_id, file_id, media_group_id)
		 VALUES ($1, $2, $3) RETURNING id;`,
		n.TextsID, n.FileID, n.MediaGroupID).Scan(&id); err != nil {
		log.Error("save a animations", logger.ErrAttr(err))
		return 0, er.New("unable to save a animations", op, err)
	}

	return id, nil
}

func (repo *pgRepository) FindByTextsID(ctx context.Context, id int) ([]string, error) {
	const op string = "animations.repository.FindByTextsID"

	var animations Animations
	if err := repo.db.QueryRow(ctx,
		`SELECT file_id, media_group_id FROM animations WHERE texts_id = $1;`,
		id).Scan(&animations.FileID, &animations.MediaGroupID); err != nil {
		return []string{}, er.New("the animations could not be found", op, err)
	}

	if animations.MediaGroupID != "" {
		rows, err := repo.db.Query(ctx,
			`SELECT file_id FROM animations WHERE media_group_id = $1;`,
			animations.MediaGroupID)
		if err != nil {
			return []string{}, er.New("the animations could not be found", op, err)
		}
		defer rows.Close()

		var animationsIDs []string
		for rows.Next() {
			var AnimationsID string
			if err := rows.Scan(&AnimationsID); err != nil {
				return nil, er.New("unable to scan data", op, err)
			}
			animationsIDs = append(animationsIDs, AnimationsID)
		}

		if err := rows.Err(); err != nil {
			return nil, er.New("error in rows: %w", op, err)
		}

		return animationsIDs, nil
	}

	return []string{animations.FileID}, nil
}

func (repo *pgRepository) UpdateByTextsID(ctx context.Context, n *Animations) error {
	const op string = "animations.repository.UpdateByID"

	if _, err := repo.db.Exec(ctx,
		`UPDATE animations
		SET file_id = $1, media_group_id = $2
		WHERE texts_id = $3;`,
		n.FileID, n.MediaGroupID, n.TextsID); err != nil {
		return er.New("unable to move animations", op, err)
	}

	return nil
}
