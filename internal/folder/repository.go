package folder

import (
	"context"
	"sync"

	"archive_bot/pkg/er"
	"archive_bot/pkg/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoFolders = er.New("there's no saved folders", "", nil)

var (
	instance *pgRepository
	once     sync.Once
)

type pgRepository struct {
	log *logger.Logger
	db  *pgxpool.Pool
}

// NewRepository creates new folder repository.
func NewRepository(ctx context.Context, log *logger.Logger, db *pgxpool.Pool) (*pgRepository, error) {
	once.Do(func() {
		instance = &pgRepository{log: log, db: db}
	})

	return instance, nil
}

// Save saves catalogue to database.
func (repo *pgRepository) Save(ctx context.Context, f *Folder) (int, error) {
	const op string = "folder.repository.Save"

	var id int
	if err := repo.db.QueryRow(ctx,
		`INSERT INTO folders (user_id, name)
		VALUES ($1, $2)
		ON CONFLICT (user_id, name) DO UPDATE 
		SET name = $2
		RETURNING id;`,
		f.UserID, f.Name).Scan(&id); err != nil {
		return 0, er.New("unable to create folder", op, err)
	}

	return id, nil
}

func (repo *pgRepository) Find(ctx context.Context, f *Folder) (string, error) {
	const op string = "folder.repository.FindOrCreate"

	var folderName string

	if err := repo.db.QueryRow(ctx,
		`SELECT name FROM folders WHERE id = $1`,
		f.ID).Scan(&folderName); err != nil {
		return "", er.New("unable to find folder(id)", op, err)
	}

	return folderName, nil
}

// TODO: maybe transaction
func (repo *pgRepository) FindOrCreate(ctx context.Context, f *Folder) (int, error) {
	const op string = "folder.repository.FindOrCreate"

	var folderID int
	if err := repo.db.QueryRow(ctx,
		`INSERT INTO folders (user_id, name)
		VALUES ($1, $2)
		ON CONFLICT (user_id, name) DO UPDATE 
		SET name = $2
		RETURNING id;`,
		f.UserID, f.Name).Scan(&folderID); err != nil {
		return 0, er.New("unable to find or create folder(id)", op, err)
	}

	return folderID, nil
}

// Save saves catalogue to database.
func (repo *pgRepository) All(ctx context.Context, f *Folder) ([]*Folder, error) {
	const op string = "folder.repository.All"

	rows, err := repo.db.Query(ctx,
		`SELECT id, name FROM folders WHERE user_id = $1;`, f.UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNoFolders
		}
		return nil, er.New("unable to get all folders", op, err)
	}
	defer rows.Close()

	catalogues := []*Folder{}
	for rows.Next() {
		var ctl Folder
		if err := rows.Scan(&ctl.ID, &ctl.Name); err != nil {
			return nil, er.New("unable to scan data", op, err)
		}
		catalogues = append(catalogues, &ctl)
	}

	// if len(catalogues) == 0 {
	// 	return nil, ErrNoFolders
	// }

	if err := rows.Err(); err != nil {
		return nil, er.New("error in rows", op, err)
	}

	return catalogues, nil
}

func (repo *pgRepository) RemoveByID(ctx context.Context, id int) error {
	const op string = "folder.repository.RemoveLast"

	if _, err := repo.db.Exec(ctx,
		`DELETE FROM folders WHERE id = ($1);`,
		id); err != nil {
		return er.New("the folder could not be removed", op, err)
	}

	return nil
}

// DefaultCatalogueID find default catalogue id.
func (repo *pgRepository) DefaultFolderID(ctx context.Context, user_id int64) (int, error) {
	const op string = "folder.repository.DefaultFolderID"

	var id int
	if err := repo.db.QueryRow(ctx,
		`SELECT id FROM folders
		WHERE user_id = $1 AND name = 'default';`,
		user_id).Scan(&id); err != nil {
		if err == pgx.ErrNoRows {
			return 0, ErrNoFolders
		}
		return 0, er.New("unable to get id of default folder", op, err)
	}

	return id, nil
}
