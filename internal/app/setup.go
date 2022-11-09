package app

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"github.com/richard-on/website/config"
	"github.com/richard-on/website/pkg/logger"
)

type App struct {
	app *fiber.App
	log logger.Logger
}

func NewApp() App {
	log := logger.NewLogger(logger.DefaultWriter(),
		config.LogLevel,
		"website-app")

	engine := html.New("./static", ".html")

	app := fiber.New(fiber.Config{
		Prefork:       false,
		ServerHeader:  "richardhere.dev",
		CaseSensitive: false,
		Views:         engine,
		ReadTimeout:   time.Second * 30,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a fiber.*Error
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Send custom error page
			err = ctx.Status(fiber.StatusNotFound).Render(strconv.Itoa(code), fiber.Map{})
			if err != nil {
				// In case the SendFile fails
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			// Return from handler
			return nil
		},
	})

	prometheus := fiberprometheus.New("richardhere.dev")
	prometheus.RegisterAt(app, "/metrics")

	app.Use(
		cors.New(cors.ConfigDefault),
		favicon.New(favicon.Config{
			File: "./static/favicon.ico",
		}),
		recover.New(),
		pprof.New(),
		prometheus.Middleware,
		logger.Middleware(
			logger.NewLogger(logger.DefaultWriter(),
				config.LogLevel,
				"website-server"), nil,
		),
	)

	app.Static("/", "./static")

	a := App{
		app: app,
		log: log,
	}

	return a
}

func (a *App) Run() {
	idleConn := make(chan struct{})

	go func() {
		// Waiting for quit signal on exit
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

		err := a.app.Shutdown()
		if err != nil {
			a.log.Fatalf(err, "could not shutdown server")
		}

		close(idleConn)
		<-quit
	}()

	if err := a.app.ListenTLS(
		":443",
		"etc/fullchain.pem",
		"etc/privkey.pem",
	); err != nil {
		a.log.Fatalf(err, "error while listening tls")
	}

	if err := a.app.Listen(":80"); err != nil {
		a.log.Fatalf(err, "error while listening at port 80")
	}

	<-idleConn
}
