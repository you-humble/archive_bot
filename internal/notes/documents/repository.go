package documents

import (
	"context"
	"sync"

	"archive_bot/pkg/er"
	"archive_bot/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoDocs = er.New("there's no saved documents", "", nil)

var (
	instance *pgRepository
	once     sync.Once
)

type pgRepository struct {
	log *logger.Logger
	db  *pgxpool.Pool
}

// NewRepository creates new documents repository.
func NewRepository(ctx context.Context, log *logger.Logger, db *pgxpool.Pool) (*pgRepository, error) {
	once.Do(func() {
		instance = &pgRepository{log: log, db: db}
	})

	return instance, nil
}

// Save saves note to database.
func (repo *pgRepository) Save(ctx context.Context, n *Docs) (int, error) {
	const op string = "documents.repository.Save"
	log := repo.log.With(logger.String("operation", op))

	var id int
	if err := repo.db.QueryRow(ctx,
		`INSERT INTO documents (texts_id, file_id, media_group_id)
		 VALUES ($1, $2, $3) RETURNING id;`,
		n.TextsID, n.FileID, n.MediaGroupID).Scan(&id); err != nil {
		log.Error("save a docs", logger.ErrAttr(err))
		return 0, er.New("unable to save a docs", op, err)
	}

	return id, nil
}

func (repo *pgRepository) FindByTextsID(ctx context.Context, id int) ([]string, error) {
	const op string = "documents.repository.FindByTextsID"

	var docs Docs
	if err := repo.db.QueryRow(ctx,
		`SELECT file_id, media_group_id FROM documents WHERE texts_id = $1;`,
		id).Scan(&docs.FileID, &docs.MediaGroupID); err != nil {
		return []string{}, er.New("the docs could not be found", op, err)
	}

	if docs.MediaGroupID != "" {
		rows, err := repo.db.Query(ctx,
			`SELECT file_id FROM documents WHERE media_group_id = $1;`,
			docs.MediaGroupID)
		if err != nil {
			return []string{}, er.New("the documents could not be found", op, err)
		}
		defer rows.Close()

		var docsIDs []string
		for rows.Next() {
			var docsID string
			if err := rows.Scan(&docsID); err != nil {
				return nil, er.New("unable to scan data", op, err)
			}
			docsIDs = append(docsIDs, docsID)
		}

		if err := rows.Err(); err != nil {
			return nil, er.New("error in rows: %w", op, err)
		}

		return docsIDs, nil
	}

	return []string{docs.FileID}, nil
}

func (repo *pgRepository) UpdateByTextsID(ctx context.Context, n *Docs) error {
	const op string = "documents.repository.UpdateByID"

	if _, err := repo.db.Exec(ctx,
		`UPDATE documents
		SET file_id = $1, media_group_id = $2
		WHERE texts_id = $3;`,
		n.FileID, n.MediaGroupID, n.TextsID); err != nil {
		return er.New("unable to move docs", op, err)
	}

	return nil
}
