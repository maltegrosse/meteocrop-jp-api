package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"os"
)

var Logger *zap.Logger = createLogger()


func createLogger() *zap.Logger {
	atom := zap.NewAtomicLevel()

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))

	defer logger.Sync()

	atom.SetLevel(zap.DebugLevel)

	return logger
}