package main

import (
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

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

	log := logger.NewLogger(config.DefaultWriter,
		zerolog.TraceLevel,
		"website-setup")

	config.GoDotEnv, err = strconv.ParseBool(os.Getenv("GODOTENV"))
	if err != nil {
		log.Infof("can't load env variable. Trying godotenv next. err: %v\n", err)
		config.GoDotEnv = true
	}

	if config.GoDotEnv {
		err = godotenv.Load()
		if err != nil {
			log.Fatal(err, "abort. Cannot load env variables using godotenv.")
		}
	}

	config.Init(log)

	if !fiber.IsChild() {
		log.Info("env and logger setup complete")
		log.Infof("richardhere.dev - version: %v; build: %v; FiberPrefork: %v; MaxCPU: %v", version, build, config.FiberPrefork, config.MaxCPU)
	}

	runtime.GOMAXPROCS(config.MaxCPU)

	err = sentry.Init(sentry.ClientOptions{
		Dsn:              config.SentryDSN,
		TracesSampleRate: config.SentryTSR,
	})
	if err != nil {
		log.Fatal(err, "sentry init failed")
	}
	defer sentry.Flush(2 * time.Second)

	if !fiber.IsChild() {
		log.Info("sentry setup complete")
	}

	server := app.NewApp()
	server.Run()

	sentry.Recover()
}
