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
var MaxCPU int
var SentryDSN string
var SentryTSR float64

func Init() error {
	var err error

	Env = os.Getenv("ENV")

	Deploy = os.Getenv("DEPLOY")

	GoDotEnv, err = strconv.ParseBool(os.Getenv("GODOTENV"))
	if err != nil {
		return err
	}

	Log = os.Getenv("LOG")

	LogCW, err = strconv.ParseBool(os.Getenv("LOG_CW"))
	if err != nil {
		return err
	}

	LogFile = os.Getenv("LOG_FILE")

	LogLevel, err = zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return err
	}

	FiberPrefork, err = strconv.ParseBool(os.Getenv("FIBER_PREFORK"))
	if err != nil {
		return err
	}

	MaxCPU, err = strconv.Atoi(os.Getenv("MAX_CPU"))
	if err != nil {
		return err
	}

	SentryDSN = os.Getenv("SENTRY_DSN")

	SentryTSR, err = strconv.ParseFloat(os.Getenv("SENTRY_TSR"), 64)
	if err != nil {
		return err
	}

	return nil
}
