package logger

// generating the logger for the application
import (
	"Go_Commerce/internal/infrastructure/config"
	"log/slog"
	"os"
)

type sloglogger struct { /// for te dependency ininversion = no direct dependent to the slog package
	logger *slog.Logger
}

func New(cfg *config.Config) Logger { // this read and return the logger based on the config
	var level slog.Level
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	options := &slog.HandlerOptions{ // only accepet the  logger that is same to lavel
		Level: level,
	}

	var h slog.Handler

	if cfg.AppEnv == "production" {
		h = slog.NewJSONHandler(os.Stdout, options) // give the json format for the production env and the text format for the dev env
	} else {
		h = slog.NewTextHandler(os.Stdout, options)
	}
	return &sloglogger{
		logger: slog.New(h), // log data is send .
	}
}

// logger level methods
func (l *sloglogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}
func (l *sloglogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *sloglogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *sloglogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

// loger interface ..
type Logger interface {
	Info(msg string, args ...any)  //  surface level information
	Warn(msg string, args ...any)  // problem
	Error(msg string, args ...any) // error /
	Debug(msg string, args ...any) //depth info
}

var _ Logger = (*sloglogger)(nil) // checking the import and throwing error if not implemented

// nopLogger discards everything. Useful in unit tests.
type nopLogger struct{}

func (nopLogger) Info(string, ...any)  {}
func (nopLogger) Warn(string, ...any)  {}
func (nopLogger) Error(string, ...any) {}
func (nopLogger) Debug(string, ...any) {}

// NewNop returns a Logger that discards all entries.
func NewNop() Logger {
	return nopLogger{}
}

/// here  i have make the logger interface
// so that that has the diffrent behaviour and
// then as per the behnaviour we define the method for each lloger  in the env life the
// and throwing the massage and then we have the struct that the struct
//  that is used ofr the wrapping the dirct depencdent so that new logger is
//  ued to read and handle the logger based oon the config

func NewLogger(cfg *config.Config) Logger {
	return New(cfg)
}
