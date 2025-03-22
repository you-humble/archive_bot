package texts

import (
	"context"
	"sync"

	"archive_bot/pkg/er"
	"archive_bot/pkg/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoTextNote = er.New("there's no saved note", "", nil)

var (
	instance *pgRepository
	once     sync.Once
)

type pgRepository struct {
	log *logger.Logger
	mu  sync.Mutex
	db  *pgxpool.Pool
}

// NewRepository creates new note repository.
func NewRepository(ctx context.Context, log *logger.Logger, db *pgxpool.Pool) (*pgRepository, error) {
	once.Do(func() {
		instance = &pgRepository{log: log, mu: sync.Mutex{}, db: db}
	})

	return instance, nil
}

// Save saves note to database.
func (repo *pgRepository) Save(ctx context.Context, n *TextNote) (int, error) {
	const op string = "texts.repository.Save"

	// TODO: one statement
	var id int
	var desc string
	if n.MediaGroupID != "" {
		repo.mu.Lock()
		if err := repo.db.QueryRow(ctx,
			`SELECT id, description FROM texts WHERE user_id = $1 AND media_group_id = $2;`,
			n.UserID, n.MediaGroupID).Scan(&id, &desc); err != nil {
			repo.log.Debug(
				"media_group_id is exists",
				logger.String("operation", op),
				logger.Int("id", id),
				logger.String("media_group_id", n.MediaGroupID),
				logger.ErrAttr(err),
			)
			if err == pgx.ErrNoRows {
				if err := repo.db.QueryRow(ctx,
					`INSERT INTO texts
						(user_id, folder_id, description, media_group_id, type)
						VALUES ($1, $2, $3, $4, $5)
						RETURNING id;`,
					n.UserID, n.FolderID, n.Description, n.MediaGroupID, n.Type,
				).Scan(&id); err != nil {
					return 0, er.New("unable to save note", op, err)
				}
			}
		} else if len(n.Description) > len(desc) {
			n.ID = id
			if err := repo.UpdateByID(ctx, n); err != nil {
				return 0, er.New("unable to set description note", op, err)
			}
		}
		repo.mu.Unlock()
	} else {
		if err := repo.db.QueryRow(ctx,
			`INSERT INTO texts
			(user_id, folder_id, description, type)
			VALUES ($1, $2, $3, $4)
			RETURNING id;`,
			n.UserID, n.FolderID, n.Description, n.Type).Scan(&id); err != nil {
			return 0, er.New("unable to save note", op, err)
		}
	}

	return id, nil
}

func (repo *pgRepository) AllFrom(ctx context.Context, n *TextNote) ([]*TextNote, error) {
	const op string = "texts.repository.AllFrom"

	rows, err := repo.db.Query(ctx,
		`SELECT id, description, type, media_group_id
		FROM texts
		WHERE user_id = $1 AND folder_id = $2
		ORDER BY created_at DESC`,
		n.UserID, n.FolderID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNoTextNote
		}
		return nil, er.New("unable to get all text notes", op, err)
	}
	defer rows.Close()

	notes := []*TextNote{}
	for rows.Next() {
		var note TextNote
		if err := rows.Scan(&note.ID, &note.Description, &note.Type, &note.MediaGroupID); err != nil {
			return nil, er.New("unable to scan data", op, err)
		}
		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		return nil, er.New("error in rows: %w", op, err)
	}

	return notes, nil
}

// Move - move  a note to catalogue.
func (repo *pgRepository) Move(ctx context.Context, n *TextNote) error {
	const op string = "texts.repository.Move"
	log := repo.log.With(logger.String("operation", op))

	if _, err := repo.db.Exec(ctx,
		`UPDATE texts SET folder_id = $1 WHERE id = $2;`,
		n.FolderID, n.ID); err != nil {
		log.Error("", logger.ErrAttr(err))
		return er.New("unable to move note", op, err)
	}

	return nil
}

func (repo *pgRepository) FindLast(ctx context.Context, userID int64) (*TextNote, error) {
	const op string = "texts.repository.FindLast"

	note := TextNote{}
	if err := repo.db.QueryRow(ctx,
		`SELECT description, created_at
		FROM texts
		WHERE user_id = $1
		ORDER BY created_at DESC 
		LIMIT 1;`,
		userID).Scan(&note.Description, &note.CreatedAt); err != nil {
		return nil, er.New("the note could not be removed", op, err)
	}

	return &note, nil
}

// Move - move  a note to catalogue.
func (repo *pgRepository) MoveLast(ctx context.Context, n *TextNote) error {
	const op string = "texts.repository.MoveLast"

	if _, err := repo.db.Exec(ctx,
		`UPDATE texts
		SET folder_id = $1
		WHERE id = (
			SELECT id FROM texts
			WHERE user_id = $2
			ORDER BY created_at DESC 
			LIMIT 1
		);`, n.FolderID, n.UserID); err != nil {
		return er.New("unable to move note", op, err)
	}

	return nil
}

// TODO: type to update
func (repo *pgRepository) UpdateByID(ctx context.Context, n *TextNote) error {
	const op string = "texts.repository.UpdateByID"

	if _, err := repo.db.Exec(ctx,
		`UPDATE texts SET description = $1 WHERE id = $2;`,
		n.Description, n.ID); err != nil {
		return er.New("unable to update note", op, err)
	}

	return nil
}

func (repo *pgRepository) RemoveByID(ctx context.Context, id int) error {
	const op string = "texts.repository.RemoveByID"

	if _, err := repo.db.Exec(ctx,
		`DELETE FROM texts WHERE id = $1;`, id); err != nil {
		return er.New("the note could not be removed", op, err)
	}

	return nil
}
