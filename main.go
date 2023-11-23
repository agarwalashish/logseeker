package main

import (
	"logseeker/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger := setupLogger()
	router := setupRouter(logger)
	http.ListenAndServe(":8080", router)
}

func setupRouter(l *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	var logsHandler handlers.LogsHandlerInterface = handlers.NewLogsHandler(l)
	r.Route("/logs", func(r chi.Router) {
		r.Post("/search", logsHandler.SearchRequest)
	})

	return r
}

func setupLogger() *zap.Logger {
	config := zap.Config{
		Encoding:         "console", // or "json"
		Level:            zap.NewAtomicLevelAt(zapcore.ErrorLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		DisableStacktrace: true,
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return logger
}
