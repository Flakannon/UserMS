package logger

import (
	"io"
	"log/slog"
	"os"
)

type LoggerInitOpts struct {
	Writer         io.Writer
	VerbosityLevel int
	Mode           string
}

func SetUpLogger(opts LoggerInitOpts) *slog.Logger {
	logLevel := new(slog.LevelVar) // This will be Info by default
	logOpts := slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}

	// Could switch on Mode field in the LoggerInitOpts
	// i.e. Prod/UAT = NewJSONHandler, Dev/Local dev = NewTextHandler
	logHandler := slog.NewJSONHandler(opts.Writer, &logOpts)

	logLevel.Set(slog.Level(opts.VerbosityLevel))

	return slog.New(logHandler)
}

// Slog does not provide a fatal by default so we can craft one in our custom package
// We can wrap the slog logger as well to achieve the same
func Fatal(err error) {
	slog.Error("Fatal error", "error", err)
	os.Exit(1)
}
