package main

import (
	"errors"
	"github.com/ispaneli/urlpresser/internal/config"
	"github.com/ispaneli/urlpresser/internal/handlers"
	"github.com/ispaneli/urlpresser/internal/logger"
	"github.com/ispaneli/urlpresser/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
)

func main() {
	cfg := config.InitConfig()

	s, err := storage.NewStorage(cfg.FileStoragePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	h := handlers.NewHandlers(s, cfg.BaseURL)

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(logger.HandlerLoggerConfig))
	e.Use(middleware.Gzip())
	e.POST("/", h.ShortingURLHandler)
	e.GET("/:id", h.RedirectingURLHandler)
	e.POST("/api/shorten", h.ShortenAPIHandler)

	if err := e.Start(cfg.HTTPAddress); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
