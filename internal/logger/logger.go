package logger

import (
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

var logger = zerolog.New(os.Stdout)

var HandlerLoggerConfig = middleware.RequestLoggerConfig{
	LogStatus: true,
	LogURI:    true,
	LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
		logger.Info().
			Time("time", v.StartTime).
			Str("uri", v.URI).
			Str("method", c.Request().Method).
			Str("latency", time.Since(v.StartTime).String()).
			Msg("request")

		logger.Info().
			Time("time", v.StartTime).
			Int("status", v.Status).
			Int64("bytes_out", c.Response().Size).
			Msg("response")

		return nil
	},
}
