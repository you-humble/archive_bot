package processor

import (
	"context"
	"strconv"
	"time"

	"archive_bot/internal/const/messages"
	"archive_bot/internal/entities"
	"archive_bot/internal/user"

	"archive_bot/pkg/logger"

	"github.com/redis/go-redis/v9"
)

type FolderService interface {
	Save(ctx context.Context, event *entities.Event) string
	RemoveByID(ctx context.Context, id int) error
	Find(ctx context.Context, event *entities.Event) (string, error)
	FindOrCreate(ctx context.Context, event *entities.Event) (int, error)
	SaveDefault(ctx context.Context, event *entities.Event) error
	All(ctx context.Context, event *entities.Event) map[string]string
	DefaultFolderID(ctx context.Context, user_id int64) int
}

type UserService interface {
	Save(ctx context.Context, event *entities.Event) error
	IsExists(ctx context.Context, event *entities.Event) error
	CountUsers(ctx context.Context) (int, error)
}

type TextNoteService interface {
	Save(ctx context.Context, event *entities.Event) (int, string)
	AllFrom(ctx context.Context, event *entities.Event) (map[int]*entities.AnswerParams, int)
	FindLast(ctx context.Context, event *entities.Event) (string, time.Time)
	Move(ctx context.Context, event *entities.Event) string
	MoveLast(ctx context.Context, event *entities.Event) string
	RemoveByID(ctx context.Context, id int) error
}

type PhotoNoteService interface {
	Save(ctx context.Context, event *entities.Event) int
	FindByTextsID(ctx context.Context, textsID int) []string
	UpdateByTextsID(ctx context.Context, event *entities.Event) error
}

type DocsNoteService interface {
	Save(ctx context.Context, event *entities.Event) int
	FindByTextsID(ctx context.Context, textsID int) []string
	UpdateByTextsID(ctx context.Context, event *entities.Event) error
}

type VideoNoteService interface {
	Save(ctx context.Context, event *entities.Event) int
	FindByTextsID(ctx context.Context, textsID int) []string
	UpdateByTextsID(ctx context.Context, event *entities.Event) error
}

type AudioNoteService interface {
	Save(ctx context.Context, event *entities.Event) int
	FindByTextsID(ctx context.Context, textsID int) []string
	UpdateByTextsID(ctx context.Context, event *entities.Event) error
}

type AniNoteService interface {
	Save(ctx context.Context, event *entities.Event) int
	FindByTextsID(ctx context.Context, textsID int) []string
	UpdateByTextsID(ctx context.Context, event *entities.Event) error
}

type VoiceNoteService interface {
	Save(ctx context.Context, event *entities.Event) int
	FindByTextsID(ctx context.Context, textsID int) []string
	UpdateByTextsID(ctx context.Context, event *entities.Event) error
}

type Storage interface {
	SetInt(ctx context.Context, key string, val int)
	Int(ctx context.Context, key string) int
	Append(ctx context.Context, key string, val int)
	PopSlice(ctx context.Context, key string) []int
}

type processor struct {
	log *logger.Logger

	user UserService

	nm noteManager
	fm folderManager

	storage Storage
}

func New(
	log *logger.Logger,
	redis *redis.Client,
	user UserService,
	folder FolderService,
	textNote TextNoteService,
	photoNote PhotoNoteService,
	docsNote DocsNoteService,
	videoNote VideoNoteService,
	audioNote AudioNoteService,
	aniNote AudioNoteService,
	voiceNote AudioNoteService,
) *processor {
	return &processor{
		log:  log,
		user: user,
		nm: newNoteManager(
			textNote, photoNote, docsNote, videoNote, audioNote, aniNote, voiceNote,
		),
		fm:      newFolderManager(folder),
		storage: newStorage(log, redis),
	}
}

func (p *processor) InitUser(ctx context.Context, event *entities.Event) {
	if err := p.user.IsExists(ctx, event); err != nil {
		if err == user.ErrUserNotExists {
			if err := p.user.Save(ctx, event); err != nil {
				p.log.Error("save user error", logger.ErrAttr(err))
			}
			if err := p.fm.service.SaveDefault(ctx, event); err != nil {
				p.log.Error("save default folder error", logger.ErrAttr(err))
			}
			return
		}
		p.log.Error("server error", logger.ErrAttr(err))
	}
}

func (p *processor) CountUsers(ctx context.Context) string {
	log := logger.L(ctx).With(logger.String("operation", "processor.CountUsers"))
	count, err := p.user.CountUsers(ctx)
	if err != nil {
		log.Error("the count of users was not found", logger.ErrAttr(err))
		return messages.Error
	}
	return "The quantity of users is " + strconv.Itoa(count)
}

const (
	msgIDPrefix       = "msg:"
	folderMsgIDPrefix = "user-msg:"
)

func (p *processor) SetInt(key string, num int) {
	p.storage.SetInt(context.Background(), key, num)
}

func (p *processor) Int(key string) int {
	return p.storage.Int(context.Background(), key)
}

func (p *processor) SetFolderMsgID(userID int64, messageID int) {
	p.storage.SetInt(
		context.Background(),
		folderMsgIDPrefix+strconv.FormatInt(userID, 10),
		messageID,
	)
}

func (p *processor) FolderMsgID(userID int64) int {
	return p.storage.Int(
		context.Background(),
		folderMsgIDPrefix+strconv.FormatInt(userID, 10),
	)
}

func (p *processor) AddMessageID(userID int64, messageID int) {
	if messageID != 0 {
		p.storage.Append(
			context.Background(),
			msgIDPrefix+strconv.FormatInt(userID, 10),
			messageID,
		)
	}
}

func (p *processor) MessageIDs(userID int64) []int {
	return p.storage.PopSlice(
		context.Background(),
		msgIDPrefix+strconv.FormatInt(userID, 10),
	)
}
