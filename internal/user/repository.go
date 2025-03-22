package user

import (
	"context"
	"sync"

	"archive_bot/pkg/er"
	"archive_bot/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound error = er.New("user not found", "", nil)

var (
	instance *pgRepository
	once     sync.Once
)

type pgRepository struct {
	log *logger.Logger
	db  *pgxpool.Pool
}

// NewRepository creates new user repository.
func NewRepository(ctx context.Context, log *logger.Logger, db *pgxpool.Pool) (*pgRepository, error) {
	once.Do(func() {
		instance = &pgRepository{log: log, db: db}
	})

	return instance, nil
}

// Save saves user to database.
func (repo *pgRepository) Save(ctx context.Context, u *User) error {
	const op string = "user.repository.Save"

	if _, err := repo.db.Exec(ctx,
		`INSERT INTO users (id, username)
		VALUES ($1, $2);`, u.ID, u.Username); err != nil {
		return er.New("unable to save user", op, err)
	}

	return nil
}

// Check that user is exists.
func (repo *pgRepository) IsExists(ctx context.Context, id int64) error {
	const op string = "user.repository.IsExists"

	var isExists bool
	if err := repo.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users
		WHERE id = $1);`, id).Scan(&isExists); err != nil {
		return er.New("unable to get user", op, err)
	}

	if !isExists {
		return ErrUserNotFound
	}

	return nil
}

// Check that user is exists.
func (repo *pgRepository) CountUsers(ctx context.Context) (int, error) {
	const op string = "user.repository.CountUsers"

	var count int
	if err := repo.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM users;`).Scan(&count); err != nil {
		return 0, er.New("unable to get count of users", op, err)
	}

	return count, nil
}
