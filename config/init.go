package config

import (
	"github.com/rs/zerolog"
	"os"
	"strconv"
)

var Env string
var Deploy string
var GoDotEnv bool
var Log string
var LogCW bool
var LogFile string
var LogLevel zerolog.Level
var FiberPrefork bool
var SentryDSN string
var SentryTSR float64

func Init() {
	var err error

	Env = os.Getenv("ENV")

	Deploy = os.Getenv("DEPLOY")

	GoDotEnv, err = strconv.ParseBool(os.Getenv("GO_DOT_ENV"))

	Log = os.Getenv("LOG")

	LogCW, err = strconv.ParseBool(os.Getenv("LOG_CW"))

	LogFile = os.Getenv("LOG_FILE")

	LogLevel, err = zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))

	FiberPrefork, err = strconv.ParseBool(os.Getenv("FIBER_PREFORK"))

	SentryDSN = os.Getenv("SENTRY_DSN")

	SentryTSR, err = strconv.ParseFloat(os.Getenv("SENTRY_TSR"), 64)

	if err != nil {
		panic(err)
	}
}
