package config

import (
	"os"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/richard-on/website/pkg/logger"
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

func Init(log logger.Logger) {
	var err error

	Env = os.Getenv("ENV")

	Deploy = os.Getenv("DEPLOY")

	GoDotEnv, err = strconv.ParseBool(os.Getenv("GODOTENV"))
	if err != nil {
		log.Infof("GODOTENV init: %v", err)
	}

	Log = os.Getenv("LOG")

	LogCW, err = strconv.ParseBool(os.Getenv("LOG_CW"))
	if err != nil {
		log.Infof("LOG_CW init: %v", err)
	}

	LogFile = os.Getenv("LOG_FILE")

	LogLevel, err = zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Infof("LOG_LEVEL init: %v", err)
	}

	FiberPrefork, err = strconv.ParseBool(os.Getenv("FIBER_PREFORK"))
	if err != nil {
		log.Infof("FIBER_PREFORK init: %v", err)
	}

	MaxCPU, err = strconv.Atoi(os.Getenv("MAX_CPU"))
	if err != nil {
		log.Infof("MAX_CPU init: %v", err)
	}

	SentryDSN = os.Getenv("SENTRY_DSN")

	SentryTSR, err = strconv.ParseFloat(os.Getenv("SENTRY_TSR"), 64)
	if err != nil {
		log.Infof("SENTRY_TSR init: %v", err)
	}
}
