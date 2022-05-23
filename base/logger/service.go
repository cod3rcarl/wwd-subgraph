package logger

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	timeLayout = "2006-01-02T15:04:05+00:00"
)

// NOTE: passing around a pointer to the logger is good, because when the log level is dynamically updated
// all users will use the updated logger

func NewLogger(cfg Config) (*zap.Logger, error) {
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format(timeLayout))
	}

	var level zapcore.Level
	if err := level.Set(cfg.Level); err != nil {
		return nil, errors.Errorf("invalid value for log level: %v", err)
	}

	atom := zap.NewAtomicLevel()
	atom.SetLevel(level)

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encCfg),
		zapcore.Lock(os.Stdout),
		atom,
	), zap.WithCaller(false))

	logger.Info("initialised logger", zap.String("level", level.String()))

	return logger, nil
}
