package main

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"

	"github.com/ispaneli/urlpresser/internal/config"
	"github.com/ispaneli/urlpresser/internal/handlers"
	"github.com/ispaneli/urlpresser/internal/storage"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.InitConfig()
	fmt.Printf("HTTP Address: %s\n", cfg.HTTPAddress)
	fmt.Printf("Base URL: %s\n", cfg.BaseURL)

	HTTPAddress := cfg.HTTPAddress
	baseURL := cfg.BaseURL

	if address := os.Getenv("SERVER_ADDRESS"); address != "" {
		HTTPAddress = address
	}
	if url := os.Getenv("BASE_URL"); url != "" {
		baseURL = url
	}

	s := storage.NewStorage()

	h := handlers.NewHandlers(s, baseURL)

	e := echo.New()
	e.Use(middleware.Logger())
	e.POST("/", h.ShortingURLHandler)
	e.GET("/:id", h.RedirectingURLHandler)

	if err := e.Start(HTTPAddress); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
