package config

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
)

var Version string
var Build string
var Env string
var LogLevelRaw string
var LogLevel zerolog.Level
var SentryDSN string

func Init() {
	var err error

	Version = os.Getenv("VERSION")
	Build = os.Getenv("BUILD")
	Env = os.Getenv("ENV")
	SentryDSN = os.Getenv("SENTRY_DSN")
	LogLevelRaw = os.Getenv("LOG_LEVEL")
	if LogLevelRaw == "" {
		_ = fmt.Errorf("abort: env variables are empty")
		panic(err)
	}

	LogLevel, err = zerolog.ParseLevel(LogLevelRaw)
	if err != nil {
		_ = fmt.Errorf("abort. Cannot parse log level. err: %v", err)
		panic(err)
	}
}
