package processor

import (
	"bytes"
	"context"
	"encoding/binary"
	"archive_bot/pkg/logger"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type storage struct {
	log *logger.Logger
	db  *redis.Client
}

func newStorage(log *logger.Logger, db *redis.Client) *storage {
	return &storage{log: log, db: db}
}

func (s *storage) SetInt(ctx context.Context, key string, val int) {
	log := s.log.With(logger.String("operation", "processor.Storage.Set"))

	if err := s.db.Set(ctx, key, val, 0).Err(); err != nil {
		log.Error(
			"failed to set value",
			logger.Int("value", val),
			logger.ErrAttr(err),
		)
	}
}

func (s *storage) Int(ctx context.Context, key string) int {
	log := s.log.With(logger.String("operation", "processor.Storage.Set"))

	val, err := s.db.Get(ctx, key).Int(); 
	if err != nil {
		log.Error(
			"failed to set value",
			logger.Int("value", val),
			logger.ErrAttr(err),
		)
		return 0
	}

	return val
}

func (s *storage) Append(ctx context.Context, key string, val int) {
	log := s.log.With(logger.String("operation", "processor.Storage.Append"))

	if err := s.db.RPush(ctx, key, val).Err(); err != nil {
		log.Error(
			"failed to set value to slice",
			logger.Int("value", val),
			logger.ErrAttr(err),
		)
	}
}

func (s *storage) PopSlice(ctx context.Context, key string) []int {
	log := s.log.With(logger.String("operation", "processor.Storage.PopSlice"))

	sl, err := redis.NewScript(`
	local list = redis.call("LRANGE", KEYS[1], 0, -1)
	redis.call("DEL", KEYS[1])
	return list
	`).Run(ctx, s.db, []string{key}).Result()
	if err != nil {
		log.Error(
			"failed to set value to slice",
			logger.String("key", key),
			logger.ErrAttr(err),
		)
	}

	rawSlice, ok := sl.([]interface{})
	if !ok {
		log.Error("wrong type", logger.String("key", key))
	}

	res := make([]int, 0, len(rawSlice))
	for _, item := range rawSlice {
		switch v := item.(type) {
		case string:
			num, err := strconv.Atoi(v)
			if err != nil {
				log.Error(
					"failed to convert to int",
					logger.String("value", v),
					logger.ErrAttr(err),
				)
			}
			res = append(res, num)
		case []byte:
			buf := bytes.NewBuffer(v)
			num, err := binary.ReadVarint(buf)
			if err != nil {
				log.Error(
					"failed to convert to int",
					logger.String("value", string(v)),
					logger.ErrAttr(err),
				)
			}
			res = append(res, int(num))
		}
	}

	return res
}
