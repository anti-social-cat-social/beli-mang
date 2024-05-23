package logger

import (
	"fmt"
	"log/slog"
)

func Info(message interface{}) {
	slog.Info(fmt.Sprintf("%v", message))
}

func Error(err error, file string) {
	slog.Error(err.Error(), slog.String("file", file))
}
