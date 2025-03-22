package voices

import (
	"context"
	"sync"

	"archive_bot/pkg/er"
	"archive_bot/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoVoice = er.New("there's no saved voices", "", nil)

var (
	instance *pgRepository
	once     sync.Once
)

type pgRepository struct {
	log *logger.Logger
	db  *pgxpool.Pool
}

// NewRepository creates new voices repository.
func NewRepository(ctx context.Context, log *logger.Logger, db *pgxpool.Pool) (*pgRepository, error) {
	once.Do(func() {
		instance = &pgRepository{log: log, db: db}
	})

	return instance, nil
}

// Save saves note to database.
func (repo *pgRepository) Save(ctx context.Context, n *Voice) (int, error) {
	const op string = "voices.repository.Save"
	log := repo.log.With(logger.String("operation", op))

	var id int
	if err := repo.db.QueryRow(ctx,
		`INSERT INTO voices (texts_id, file_id)
		 VALUES ($1, $2) RETURNING id;`,
		n.TextsID, n.FileID).Scan(&id); err != nil {
		log.Error("save a Voice", logger.ErrAttr(err))
		return 0, er.New("unable to save a voice", op, err)
	}

	return id, nil
}

func (repo *pgRepository) FindByTextsID(ctx context.Context, id int) ([]string, error) {
	const op string = "voices.repository.FindByTextsID"

	var voice Voice
	if err := repo.db.QueryRow(ctx,
		`SELECT file_id FROM voices WHERE texts_id = $1;`,
		id).Scan(&voice.FileID); err != nil {
		return []string{}, er.New("the voice could not be found", op, err)
	}

	return []string{voice.FileID}, nil
}

func (repo *pgRepository) UpdateByTextsID(ctx context.Context, n *Voice) error {
	const op string = "voices.repository.UpdateByID"

	if _, err := repo.db.Exec(ctx,
		`UPDATE voices SET file_id = $1 WHERE texts_id = $2;`,
		n.FileID, n.TextsID); err != nil {
		return er.New("unable to move voice", op, err)
	}

	return nil
}
