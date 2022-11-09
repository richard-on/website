package logger

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/richard-on/website/config"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	log zerolog.Logger
}

var PrettyPrint = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC1123}

func DefaultWriter() io.Writer {
	switch config.Env {
	case "win-dev":
		return PrettyPrint

	case "linux-dev":
		return os.Stdout

	case "linux-dev-file":
		logFile, err := os.Create("logs/website.log")
		if err != nil {
			sentry.CaptureException(err)
			return os.Stdout
		}
		return logFile

	default:
		return os.Stdout
	}
}

func NewLogger(out io.Writer, level zerolog.Level, service string) Logger {

	return Logger{log: zerolog.New(out).
		Level(level).
		With().
		Timestamp().
		Int("pid", os.Getpid()).
		Str("service", service).
		Logger(),
	}
}

func (l Logger) Println(v ...interface{}) {
	l.log.Print(fmt.Sprint(v...))
}

func (l Logger) Printf(format string, v ...interface{}) {
	l.log.Printf(format, v...)
}

func (l Logger) Debug(i ...interface{}) {
	l.log.Debug().Msgf(fmt.Sprint(i...))
}

func (l Logger) Dedugf(format string, i ...interface{}) {
	l.log.Debug().Msgf(format, i...)
}

func (l Logger) Info(i ...interface{}) {
	l.log.Info().Msgf(fmt.Sprint(i...))
}

func (l Logger) Infof(format string, i ...interface{}) {
	l.log.Info().Msgf(format, i...)
}

func (l Logger) Error(err error, msg string) {
	l.log.Error().Err(err).Msg(msg)
	sentry.CaptureException(err)
}

func (l Logger) Errorf(err error, format string, i ...interface{}) {
	l.log.Error().Err(err).Msgf(format, i...)
	sentry.CaptureException(err)
}

func (l Logger) Fatal(err error, msg string) {
	l.log.Fatal().Err(err).Msg(msg)
	sentry.CaptureException(err)
}

func (l Logger) Fatalf(err error, format string, i ...interface{}) {
	l.log.Fatal().Err(err).Msgf(format, i...)
	sentry.CaptureException(err)
}
