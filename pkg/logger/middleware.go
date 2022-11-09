package logger

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type fiberLog struct {
	RID        string
	RemoteIP   string
	Host       string
	Method     string
	Path       string
	Protocol   string
	StatusCode int
	Latency    float64
	Error      error
	Stack      []byte
}

func (f *fiberLog) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("id", f.RID).
		Str("remote_ip", f.RemoteIP).
		Str("host", f.Host).
		Str("method", f.Method).
		Str("path", f.Path).
		Str("protocol", f.Protocol).
		Int("status_code", f.StatusCode).
		Float64("latency", f.Latency).
		Str("tag", "request")

	if f.Error != nil {
		e.Err(f.Error)
	}

	if f.Stack != nil {
		e.Bytes("stack", f.Stack)
	}
}

func Middleware(logger Logger, filter func(ctx *fiber.Ctx) bool) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		if filter != nil && filter(ctx) {
			ctx.Next()
			return nil
		}

		start := time.Now()

		rid := ctx.Get(fiber.HeaderXRequestID)
		if rid == "" {
			rid = uuid.New().String()
			ctx.Set(fiber.HeaderXRequestID, rid)
		}

		event := &fiberLog{
			RID:      rid,
			RemoteIP: ctx.IP(),
			Host:     ctx.Hostname(),
			Method:   ctx.Method(),
			Path:     ctx.Path(),
			Protocol: ctx.Protocol(),
		}

		defer func() {
			rec := recover()

			if rec != nil {
				err, ok := rec.(error)
				if !ok {
					err = fmt.Errorf("%v", rec)
				}

				event.Error = err
				event.Stack = debug.Stack()

				ctx.Status(fiber.StatusInternalServerError)
				ctx.JSON(map[string]interface{}{
					"status": http.StatusText(fiber.StatusInternalServerError),
				})
			}

			event.StatusCode = ctx.Response().StatusCode()
			event.Latency = time.Since(start).Seconds()

			switch {
			case rec != nil:
				logger.log.Error().EmbedObject(event).Msg("server_fatalpanic")
			case event.StatusCode >= 500:
				logger.log.Error().EmbedObject(event).Msg("server_error")
			case event.StatusCode >= 400:
				logger.log.Info().EmbedObject(event).Msg("client_error")
			case event.StatusCode >= 300:
				logger.log.Info().EmbedObject(event).Msg("redirect")
			case event.StatusCode >= 200:
				logger.log.Info().EmbedObject(event).Msg("success")
			case event.StatusCode >= 100:
				logger.log.Info().EmbedObject(event).Msg("informative")
			default:
				logger.log.Warn().EmbedObject(event).Msg("status_unknown")
			}
		}()

		err := ctx.Next()
		if err != nil {
			return err
		}

		return nil
	}
}
