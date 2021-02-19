package logger

import "go.uber.org/zap/zapcore"

type LogLevelStr string

const (
	LevelDebug LogLevelStr = "debug"
	LevelFatal LogLevelStr = "fatal"
	LevelError LogLevelStr = "error"
	LevelWarn  LogLevelStr = "warn"
	LevelInfo  LogLevelStr = "info"
)

func Str2ZapLevel(level LogLevelStr) (zapcore.Level, error) {
	l := zapcore.DebugLevel
	return l, l.Set(string(level))
}
