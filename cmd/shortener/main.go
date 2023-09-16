package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/ispaneli/urlpresser/cmd/shortener/config"
)

// ShortURLMap stores the mapping of shortened URLs to original URLs.
var ShortURLMap = make(map[string]string)

// OriginalURLMap stores the mapping of original URLs to shortened URLs.
var OriginalURLMap = make(map[string]string)

func generateShortURL() string {
	/*
	   Generating a random string.

	   Description:
	   Function to generate a unique identifier for the shortened link.
	*/
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 8)
	for i := range result {
		result[i] = chars[randomInt(0, len(chars)-1)]
	}
	return string(result)
}

func randomInt(min int, max int) int {
	/*
	   Generating a random number.
	*/
	return min + rand.Intn(max-min+1)
}

func shortingURLHandler(c *fiber.Ctx, baseURL string) error {
	/*
	   Handler for URL shortening.
	*/
	originalURL := strings.TrimSpace(string(c.Body()))
	if originalURL == "" {
		return c.Status(http.StatusBadRequest).SendString("Invalid request")
	}

	// Check if the original URL already exists in OriginalURLMap.
	existShortURL, ok := OriginalURLMap[originalURL]
	if ok {
		responseURL := fmt.Sprintf(`%s%s`, baseURL, existShortURL)
		return c.Status(http.StatusCreated).SendString(responseURL)
	}

	// Generate a unique identifier for the shortened link.
	for {
		shortURL := generateShortURL()
		if _, ok := ShortURLMap[shortURL]; !ok {
			ShortURLMap[shortURL] = originalURL
			OriginalURLMap[originalURL] = shortURL

			responseURL := fmt.Sprintf(`%s%s`, baseURL, shortURL)
			return c.Status(http.StatusCreated).SendString(responseURL)
		}
	}
}

func redirectingURLHandler(c *fiber.Ctx) error {
	/*
	   Handler for redirecting to the original URL.
	*/
	shortURL := c.Params("id")
	originalURL, ok := ShortURLMap[shortURL]
	if !ok {
		return c.Status(http.StatusBadRequest).SendString("Invalid short URL")
	}

	c.Set("Location", originalURL)
	c.Status(http.StatusTemporaryRedirect)
	return nil
}

func main() {
	// Parsing command args.
	cfg := config.InitConfig()
	fmt.Printf("HTTP Address: %s\n", cfg.HTTPAddress)
	fmt.Printf("Base URL: %s\n", cfg.BaseURL)

	HTTPAddress := cfg.HTTPAddress
	baseURL := cfg.BaseURL

	// Parsing environment variables.
	if address := os.Getenv("SERVER_ADDRESS"); address != "" {
		HTTPAddress = address
	}
	if url := os.Getenv("BASE_URL"); url != "" {
		baseURL = url
	}

	// Creating and configuring web-application.
	app := fiber.New()

	app.Post(`/`, func(c *fiber.Ctx) error {
		return shortingURLHandler(c, baseURL)
	})
	app.Get(`/:id`, redirectingURLHandler)

	err := app.Listen(HTTPAddress)
	if err != nil {
		log.Fatal(err)
	}
}
