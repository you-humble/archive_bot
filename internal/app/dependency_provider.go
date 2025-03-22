package app

import (
	"context"

	"archive_bot/internal/config"
	"archive_bot/internal/folder"
	"archive_bot/internal/notes/animations"
	"archive_bot/internal/notes/audios"
	"archive_bot/internal/notes/documents"
	"archive_bot/internal/notes/photos"
	"archive_bot/internal/notes/texts"
	"archive_bot/internal/notes/videos"
	"archive_bot/internal/notes/voices"
	"archive_bot/internal/processor"
	"archive_bot/internal/router"
	"archive_bot/internal/user"

	"archive_bot/pkg/closer"
	"archive_bot/pkg/database/postgres"
	storage "archive_bot/pkg/database/redis"
	"archive_bot/pkg/er"
	"archive_bot/pkg/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Router interface {
	RouteMessage(ctx context.Context, b *bot.Bot, update *models.Update)
	RouteCallbackQuery(ctx context.Context, b *bot.Bot, update *models.Update)
	RouteAdminMessage(ctx context.Context, b *bot.Bot, update *models.Update)
	RouteAdminCallback(ctx context.Context, b *bot.Bot, update *models.Update)
}

type dependencyProvider struct {
	config *config.Config
	logger *logger.Logger

	redis *redis.Client

	db               *pgxpool.Pool
	userRepository   user.Repository
	folderRepository folder.Repository
	textRepository   texts.Repository
	photoRepository  photos.Repository
	docsRepository   documents.Repository
	videoRepository  videos.Repository
	audioRepository  audios.Repository
	aniRepository    animations.Repository
	voiceRepository  voices.Repository

	userService   processor.UserService
	folderService processor.FolderService
	textService   processor.TextNoteService
	photoService  processor.PhotoNoteService
	docsService   processor.DocsNoteService
	videoService  processor.VideoNoteService
	audioService  processor.VideoNoteService
	aniService    processor.AniNoteService
	voiceService  processor.VoiceNoteService

	processor router.Processor

	router Router
}

func newDependencyProvider() *dependencyProvider {
	return &dependencyProvider{}
}

func (dp *dependencyProvider) Config() *config.Config {
	const op = "app.Config"

	if dp.config == nil {
		cfg, err := config.New()
		if err != nil {
			panic(er.New("failed to get config", op, err))
		}

		dp.config = cfg
	}

	return dp.config
}

// TODO: log level from config or flag
func (dp *dependencyProvider) Logger() *logger.Logger {
	if dp.logger == nil {
		dp.logger = logger.NewLogger(
			logger.WithLevel(dp.Config().LogLevel),
		)
	}

	return dp.logger
}

func (dp *dependencyProvider) Redis(ctx context.Context) *redis.Client {
	const op = "app.Redis"

	if dp.redis == nil {
		db, err := storage.NewConnect(ctx, dp.Config().Redis.Options())
		if err != nil {
			panic(er.New("failed to connect to redis", op, err))
		}
		closer.Add(db.Close)
		dp.Logger().Debug("✓ connected to redis")

		dp.redis = db
	}

	return dp.redis
}

func (dp *dependencyProvider) DB(ctx context.Context) *pgxpool.Pool {
	const op = "app.DB"

	if dp.db == nil {
		db, err := postgres.NewConnect(ctx, dp.Config().PostgresURL)
		if err != nil {
			panic(er.New("failed to connect to the database", op, err))
		}
		closer.Add(func() error {
			db.Close()
			return nil
		})
		dp.Logger().Debug("✓ connected to folder_holder_db db")

		dp.db = db
	}

	return dp.db
}

func (dp *dependencyProvider) UserRepository(ctx context.Context) user.Repository {
	const op = "app.UserRepository"

	if dp.userRepository == nil {
		repo, err := user.NewRepository(ctx, dp.Logger(), dp.DB(ctx))
		if err != nil {
			panic(er.New("failed to create user repository", op, err))
		}

		dp.userRepository = repo
	}

	return dp.userRepository
}

func (dp *dependencyProvider) FolderRepository(ctx context.Context) folder.Repository {
	const op = "app.FolderRepository"

	if dp.folderRepository == nil {
		repo, err := folder.NewRepository(ctx, dp.Logger(), dp.DB(ctx))
		if err != nil {
			panic(er.New("failed to create folder repository", op, err))
		}

		dp.folderRepository = repo
	}

	return dp.folderRepository
}

func (dp *dependencyProvider) TextNoteRepository(ctx context.Context) texts.Repository {
	const op = "app.TextNoteRepository"

	if dp.textRepository == nil {
		repo, err := texts.NewRepository(ctx, dp.Logger(), dp.DB(ctx))
		if err != nil {
			panic(er.New("failed to create text note repository", op, err))
		}

		dp.textRepository = repo
	}

	return dp.textRepository
}

func (dp *dependencyProvider) PhotoNoteRepository(ctx context.Context) photos.Repository {
	const op = "app.PhotoNoteRepository"

	if dp.photoRepository == nil {
		repo, err := photos.NewRepository(ctx, dp.Logger(), dp.DB(ctx))
		if err != nil {
			panic(er.New("failed to create photo note repository", op, err))
		}

		dp.photoRepository = repo
	}

	return dp.photoRepository
}

func (dp *dependencyProvider) DocsNoteRepository(ctx context.Context) documents.Repository {
	const op = "app.DocsNoteRepository"

	if dp.docsRepository == nil {
		repo, err := documents.NewRepository(ctx, dp.Logger(), dp.DB(ctx))
		if err != nil {
			panic(er.New("failed to create documents note repository", op, err))
		}

		dp.docsRepository = repo
	}

	return dp.docsRepository
}

func (dp *dependencyProvider) VideoNoteRepository(ctx context.Context) videos.Repository {
	const op = "app.VideoNoteRepository"

	if dp.videoRepository == nil {
		repo, err := videos.NewRepository(ctx, dp.Logger(), dp.DB(ctx))
		if err != nil {
			panic(er.New("failed to create video note repository", op, err))
		}

		dp.videoRepository = repo
	}

	return dp.videoRepository
}

func (dp *dependencyProvider) AudioNoteRepository(ctx context.Context) audios.Repository {
	const op = "app.AudioNoteRepository"

	if dp.audioRepository == nil {
		repo, err := audios.NewRepository(ctx, dp.Logger(), dp.DB(ctx))
		if err != nil {
			panic(er.New("failed to create audio note repository", op, err))
		}

		dp.audioRepository = repo
	}

	return dp.audioRepository
}

func (dp *dependencyProvider) AniNoteRepository(ctx context.Context) animations.Repository {
	const op = "app.AniNoteRepository"

	if dp.aniRepository == nil {
		repo, err := animations.NewRepository(ctx, dp.Logger(), dp.DB(ctx))
		if err != nil {
			panic(er.New("failed to create animation note repository", op, err))
		}

		dp.aniRepository = repo
	}

	return dp.aniRepository
}

func (dp *dependencyProvider) VoiceNoteRepository(ctx context.Context) voices.Repository {
	const op = "app.VoicesNoteRepository"

	if dp.voiceRepository == nil {
		repo, err := voices.NewRepository(ctx, dp.Logger(), dp.DB(ctx))
		if err != nil {
			panic(er.New("failed to create voice note repository", op, err))
		}

		dp.voiceRepository = repo
	}

	return dp.voiceRepository
}

func (dp *dependencyProvider) UserService(ctx context.Context) processor.UserService {
	if dp.userService == nil {
		dp.userService = user.NewService(ctx, dp.Logger(), dp.UserRepository(ctx))
	}

	return dp.userService
}

func (dp *dependencyProvider) FolderService(ctx context.Context) processor.FolderService {
	if dp.folderService == nil {
		dp.folderService = folder.NewService(ctx, dp.Logger(), dp.FolderRepository(ctx))
	}

	return dp.folderService
}

func (dp *dependencyProvider) TextNoteService(ctx context.Context) processor.TextNoteService {
	if dp.textService == nil {
		dp.textService = texts.NewService(ctx, dp.Logger(), dp.TextNoteRepository(ctx))
	}

	return dp.textService
}

func (dp *dependencyProvider) PhotoNoteService(ctx context.Context) processor.PhotoNoteService {
	if dp.photoService == nil {
		dp.photoService = photos.NewService(ctx, dp.Logger(), dp.PhotoNoteRepository(ctx))
	}

	return dp.photoService
}

func (dp *dependencyProvider) DocsNoteService(ctx context.Context) processor.DocsNoteService {
	if dp.docsService == nil {
		dp.docsService = documents.NewService(ctx, dp.Logger(), dp.DocsNoteRepository(ctx))
	}

	return dp.docsService
}

func (dp *dependencyProvider) VideoNoteService(ctx context.Context) processor.VideoNoteService {
	if dp.videoService == nil {
		dp.videoService = videos.NewService(ctx, dp.Logger(), dp.VideoNoteRepository(ctx))
	}

	return dp.videoService
}

func (dp *dependencyProvider) AudioNoteService(ctx context.Context) processor.AudioNoteService {
	if dp.audioService == nil {
		dp.audioService = audios.NewService(ctx, dp.Logger(), dp.AudioNoteRepository(ctx))
	}

	return dp.audioService
}
func (dp *dependencyProvider) AniNoteService(ctx context.Context) processor.AniNoteService {
	if dp.aniService == nil {
		dp.aniService = animations.NewService(ctx, dp.Logger(), dp.AniNoteRepository(ctx))
	}

	return dp.aniService
}
func (dp *dependencyProvider) VoiceNoteService(ctx context.Context) processor.VoiceNoteService {
	if dp.voiceService == nil {
		dp.voiceService = voices.NewService(ctx, dp.Logger(), dp.VoiceNoteRepository(ctx))
	}

	return dp.voiceService
}

func (dp *dependencyProvider) Processor(ctx context.Context) router.Processor {
	if dp.processor == nil {
		dp.processor = processor.New(
			dp.Logger(),
			dp.Redis(ctx),
			dp.UserService(ctx),
			dp.FolderService(ctx),
			dp.TextNoteService(ctx),
			dp.PhotoNoteService(ctx),
			dp.DocsNoteService(ctx),
			dp.VideoNoteService(ctx),
			dp.AudioNoteService(ctx),
			dp.AniNoteService(ctx),
			dp.VoiceNoteService(ctx),
		)
	}

	return dp.processor
}

func (dp *dependencyProvider) Router(ctx context.Context) Router {
	if dp.router == nil {
		dp.router = router.New(
			dp.Logger(),
			dp.Config().AdminID,
			dp.Processor(ctx),
		)
	}

	return dp.router
}
