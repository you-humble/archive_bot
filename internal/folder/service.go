package folder

import (
	"context"
	"strconv"

	"archive_bot/internal/const/buttons"
	"archive_bot/internal/const/messages"
	"archive_bot/internal/entities"
	"archive_bot/pkg/logger"
)

type Repository interface {
	Save(ctx context.Context, c *Folder) (int, error)
	Find(ctx context.Context, f *Folder) (string, error)
	FindOrCreate(ctx context.Context, c *Folder) (int, error)
	All(ctx context.Context, c *Folder) ([]*Folder, error)
	RemoveByID(ctx context.Context, id int) error
	DefaultFolderID(ctx context.Context, user_id int64) (int, error)
}

type service struct {
	log             *logger.Logger
	repo            Repository
	defaultFolderID map[int64]int
}

func NewService(ctx context.Context, log *logger.Logger, repo Repository) *service {
	return &service{log: log, repo: repo, defaultFolderID: make(map[int64]int)}
}

func (s *service) Save(ctx context.Context, event *entities.Event) string {
	log := s.log.With(logger.String("operation", "folder.service.Save"))

	if _, err := s.repo.Save(ctx, &Folder{
		UserID: event.Meta.UserID, Name: event.Text,
	}); err != nil {
		log.Error("failed to save folder", logger.ErrAttr(err))
		return messages.Error
	}

	return messages.FolderCreated
}

func (s *service) RemoveByID(ctx context.Context, id int) error {
	return s.repo.RemoveByID(ctx, id)
}

func (s *service) FindOrCreate(ctx context.Context, event *entities.Event) (int, error) {
	return s.repo.FindOrCreate(ctx, &Folder{
		UserID: event.Meta.UserID, Name: event.Text,
	})
}

func (s *service) Find(ctx context.Context, event *entities.Event) (string, error) {
	return s.repo.Find(ctx, &Folder{ID: event.FolderID})
}

func (s *service) SaveDefault(ctx context.Context, event *entities.Event) error {
	log := s.log.With(logger.String("operation", "folder.service.Save"))

	c := &Folder{UserID: event.Meta.UserID, Name: "default"}
	FolderID, err := s.repo.Save(ctx, c)
	if err != nil {
		log.Error(
			"failed to save Folder",
			logger.ErrAttr(err),
		)
		return err
	}

	s.defaultFolderID[event.Meta.UserID] = FolderID
	return nil
}

func (s *service) DefaultFolderID(ctx context.Context, user_id int64) int {
	defaultFolderID, ok := s.defaultFolderID[user_id]
	if !ok {
		defaultFolderID, err := s.repo.DefaultFolderID(ctx, user_id)
		if err != nil {
			return 0
		}
		s.defaultFolderID[user_id] = defaultFolderID
	}
	return defaultFolderID
}

func (s *service) All(ctx context.Context, event *entities.Event) map[string]string {
	log := s.log.With(logger.String("operation", "Folder.service.AllAsInlineButtons"))

	folders, err := s.repo.All(ctx, &Folder{UserID: event.Meta.UserID})
	if err != nil {
		if err == ErrNoFolders {
			log.Info("Folders is empty", logger.ErrAttr(err))
			return nil
		}
		log.Error("failed to get all Folders", logger.ErrAttr(err))
		return nil
	}

	res := make(map[string]string, len(folders))
	for _, f := range folders {
		if f.Name == "default" {
			res[buttons.Prefix+strconv.Itoa(f.ID)] = buttons.DefaultFolderName
			continue
		}
		res[buttons.Prefix+strconv.Itoa(f.ID)] = f.Name
	}

	return res
}
