package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"

	"github.com/richard-on/website/config"
)

type Logger struct {
	log zerolog.Logger
}

var DefaultWriter = NewWriter()

func NewWriter() io.Writer {

	var out io.Writer

	switch config.Log {
	case "stdout":
		out = os.Stdout
	case "stderr":
		out = os.Stderr
	case "file":
		file, err := os.OpenFile("logs/"+config.LogFile, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}

		out = file

	default:
		out = os.Stdout
	}

	if config.LogCW {
		return zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC1123}
	}

	return out
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
