package main

import (
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/richard-on/website/config"
	"github.com/richard-on/website/internal/app"
	"github.com/richard-on/website/pkg/logger"
)

func main() {

	if len(os.Args) > 1 {
		if os.Args[1] == "win-dev" {
			config.Env = "win-dev"

			err := godotenv.Load()
			if err != nil {
				_ = fmt.Errorf("abort. Cannot load env variables. err: %v", err)
				panic(err)
			}
		}
	}

	config.Init()

	log := logger.NewLogger(logger.DefaultWriter(),
		config.LogLevel,
		"website-setup")

	log.Info("env setup complete")

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.SentryDSN,
		TracesSampleRate: 1.0,
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
