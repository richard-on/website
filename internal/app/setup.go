package app

import (
	"strconv"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	log := logger.NewLogger(config.DefaultWriter,
		config.LogLevel,
		"website-app")

	engine := html.New("./public", ".html")

	app := fiber.New(fiber.Config{
		Prefork:       config.FiberPrefork,
		ServerHeader:  "richardhere.dev",
		CaseSensitive: false,
		Views:         engine,
		ReadTimeout:   time.Second * 30,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {

			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			err = ctx.Status(fiber.StatusNotFound).Render(strconv.Itoa(code), fiber.Map{})
			if err != nil {
				// In case the SendFile fails
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			return nil
		},
	})

	prometheus := fiberprometheus.New("richardhere.dev")
	prometheus.RegisterAt(app, "/metrics")

	app.Use(
		cors.New(cors.ConfigDefault),
		recover.New(),
		pprof.New(
			pprof.Config{Next: func(c *fiber.Ctx) bool {
				return config.Env != "dev"
			}}),
		prometheus.Middleware,
		logger.Middleware(
			logger.NewLogger(config.DefaultWriter,
				config.LogLevel,
				"website-server"), nil,
		),
	)

	app.Static("/", "./public")

	return App{
		app: app,
		log: log,
	}
}
