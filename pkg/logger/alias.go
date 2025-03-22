package logger

import (
	"log/slog"
	"time"
)

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

type (
	Logger         = slog.Logger
	Attr           = slog.Attr
	Level          = slog.Level
	Handler        = slog.Handler
	Value          = slog.Value
	HandlerOptions = slog.HandlerOptions
	LogValuer      = slog.LogValuer
)

var (
	NewTextHandler = slog.NewTextHandler
	NewJSONHandler = slog.NewJSONHandler
	New            = slog.New
	SetDefault     = slog.SetDefault

	String   = slog.String
	Bool     = slog.Bool
	Float64  = slog.Float64
	Any      = slog.Any
	Duration = slog.Duration
	Int      = slog.Int
	Int64    = slog.Int64
	Uint64   = slog.Uint64

	GroupValue = slog.GroupValue
	Group      = slog.Group

	Default = slog.Default
)

func ErrAttr(err error) Attr {
	return slog.String("error", err.Error())
}

func Float32Attr(key string, val float32) Attr {
	return slog.Float64(key, float64(val))
}

func UInt32Attr(key string, val uint32) Attr {
	return slog.Int(key, int(val))
}

func Int32Attr(key string, val int32) Attr {
	return slog.Int(key, int(val))
}

func TimeAttr(key string, time time.Time) Attr {
	return slog.String(key, time.String())
}
