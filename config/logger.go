package config

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"time"
)

var DefaultWriter = NewWriter()

func NewWriter() io.Writer {

	var out io.Writer

	switch Log {
	case "stdout":
		out = os.Stdout
	case "stderr":
		out = os.Stderr
	case "file":
		file, err := os.OpenFile("logs/"+LogFile, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}

		out = file

	default:
		out = os.Stdout
	}

	if LogCW {
		return zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC1123}
	}

	return out
}
