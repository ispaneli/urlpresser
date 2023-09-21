package main

import (
	"fmt"
	"github.com/ispaneli/urlpresser/cmd/shortener/config"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/valyala/fasthttp"
)

type URLMap struct {
	// ShortURLMap stores the mapping of shortened URLs to original URLs.
	ShortURLMap map[string]string
	// OriginalURLMap stores the mapping of original URLs to shortened URLs.
	OriginalURLMap map[string]string

	sync.Mutex
}

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

func shortingURLHandler(ctx *fasthttp.RequestCtx, baseURL string) {
	/*
	   Handler for URL shortening.
	*/
	originalURL := strings.TrimSpace(string(ctx.PostBody()))
	if originalURL == "" {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("Invalid request")
		return
	}

	// Check if the original URL already exists in OriginalURLMap.
	existShortURL, ok := OriginalURLMap[originalURL]
	if ok {
		responseURL := fmt.Sprintf(`%s%s`, baseURL, existShortURL)
		ctx.SetStatusCode(http.StatusCreated)
		ctx.SetBodyString(responseURL)
		return
	}

	// Generate a unique identifier for the shortened link.
	for {
		shortURL := generateShortURL()
		if _, ok := ShortURLMap[shortURL]; !ok {
			ShortURLMap[shortURL] = originalURL
			OriginalURLMap[originalURL] = shortURL

			responseURL := fmt.Sprintf(`%s%s`, baseURL, shortURL)
			ctx.SetStatusCode(http.StatusCreated)
			ctx.SetBodyString(responseURL)
			return
		}
	}
}

func redirectingURLHandler(ctx *fasthttp.RequestCtx) {
	/*
	   Handler for redirecting to the original URL.
	*/
	shortURL := ctx.UserValue("id").(string)
	originalURL, ok := ShortURLMap[shortURL]
	if !ok {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("Invalid short URL")
		return
	}

	ctx.Redirect(originalURL, http.StatusTemporaryRedirect)
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
	server := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			switch string(ctx.Path()) {
			case "/":
				if ctx.IsPost() {
					shortingURLHandler(ctx, baseURL)
				} else {
					ctx.SetStatusCode(http.StatusMethodNotAllowed)
				}
			default:
				redirectingURLHandler(ctx)
			}
		},
	}

	err := server.ListenAndServe(HTTPAddress)
	if err != nil {
		log.Fatal(err)
	}
}
