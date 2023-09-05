package logger

import (
	"go-graphql-mongo-server/config"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

func Initialize() {

	loggerEncoder := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller",
		FunctionKey:      "func",
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.LowercaseLevelEncoder,
		EncodeTime:       zapcore.RFC3339TimeEncoder,
		EncodeDuration:   zapcore.StringDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		ConsoleSeparator: "|\t|",
	}

	var loggerCore zapcore.Core

	if config.Store.ProductionMode {

		// For production, output logs in JSON Format
		loggerCore = zapcore.NewCore(
			zapcore.NewJSONEncoder(loggerEncoder),
			zapcore.AddSync(os.Stdout),
			zapcore.InfoLevel,
		)

	} else {

		// For local dev mode, output logs in Console Format
		loggerEncoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
		loggerCore = zapcore.NewCore(
			zapcore.NewConsoleEncoder(loggerEncoder),
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		)

	}

	Log = zap.New(loggerCore, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)).Sugar()

}
