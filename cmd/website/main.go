package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/richard-on/website/config"
	"github.com/richard-on/website/internal/app"
	"github.com/richard-on/website/pkg/logger"
)

var (
	version string
	build   string
)

func main() {
	var err error

	config.GoDotEnv, err = strconv.ParseBool(os.Getenv("GODOTENV"))
	if err != nil {
		fmt.Printf("can't load env variable. Trying godotenv next. err: %v\n", err)
		config.GoDotEnv = true
	}

	if config.GoDotEnv {
		err = godotenv.Load()
		if err != nil {
			_ = fmt.Errorf("abort. Cannot load env variables using godotenv. err: %v", err)
			panic(err)
		}
	}

	config.Init()

	log := logger.NewLogger(logger.DefaultWriter,
		config.LogLevel,
		"website-setup")

	log.Info("env and logger setup complete")
	log.Infof("richardhere.dev - version: %v; build: %v", version, build)

	err = sentry.Init(sentry.ClientOptions{
		Dsn:              config.SentryDSN,
		TracesSampleRate: config.SentryTSR,
	})
	if err != nil {
		log.Fatal(err, "sentry init failed")
	}
	defer sentry.Flush(2 * time.Second)

	log.Info("sentry setup complete")

	server := app.NewApp()
	server.Run()

	sentry.Recover()
}
