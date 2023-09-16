package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestShortingURLRoute(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		body         string
		contentType  string
		expectedCode int
	}{
		{
			description:  "Get HTTP status 201",
			route:        `/`,
			body:         "https://practicum.yandex.ru/",
			contentType:  "text/plain",
			expectedCode: http.StatusCreated,
		},
		{
			description:  "Get HTTP status 400 (body was empty)",
			route:        `/`,
			body:         "",
			contentType:  "text/plain",
			expectedCode: http.StatusBadRequest,
		},
		{
			description:  "Get HTTP status 201 (existing URL)",
			route:        `/`,
			body:         "https://practicum.yandex.ru/",
			contentType:  "text/plain",
			expectedCode: http.StatusCreated,
		},
	}

	app := fiber.New()
	app.Post(`/`, func(c *fiber.Ctx) error {
		return shortingURLHandler(c, "http://localhost:8080/")
	})

	var originalURLMap = make(map[string]string)

	for _, test := range tests {
		// Create a POST request with the test data.
		req := httptest.NewRequest(http.MethodPost, `/`, bytes.NewBufferString(test.body))
		req.Header.Set("Content-Type", test.contentType)
		resp, _ := app.Test(req, 1)

		// Assert that the response status code matches the expected code.
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		// Read the response body to get the short URL.
		bodyBuffer := new(bytes.Buffer)
		bodyBuffer.ReadFrom(resp.Body)
		shortURL := bodyBuffer.String()

		// Check if the original URL already exists in the map.
		if existShortURL, ok := originalURLMap[test.body]; ok {
			// If it exists, assert that the short URL matches the existing one.
			assert.Equalf(t, existShortURL, shortURL, test.description)
		} else {
			// If it doesn't exist, add it to the map.
			originalURLMap[test.body] = shortURL
		}
	}
}

func TestRedirectingURLRoute(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		body         string
		contentType  string
		expectedCode int
	}{
		{
			description:  "Get HTTP status 200",
			route:        `/`,
			body:         "https://practicum.yandex.ru/",
			contentType:  "text/plain",
			expectedCode: http.StatusCreated,
		},
		{
			description:  "Get HTTP status 201 (existing URL)",
			route:        `/`,
			body:         "https://practicum.yandex.ru/",
			contentType:  "text/plain",
			expectedCode: http.StatusCreated,
		},
	}

	app := fiber.New()
	app.Post(`/`, func(c *fiber.Ctx) error {
		return shortingURLHandler(c, "http://localhost:8080/")
	})
	app.Get(`/:id`, redirectingURLHandler)

	for _, test := range tests {
		// Create a POST request to shorten the URL.
		postReq := httptest.NewRequest(http.MethodPost, `/`, bytes.NewBufferString(test.body))
		postReq.Header.Set("Content-Type", test.contentType)
		postResp, _ := app.Test(postReq, 1)

		// Assert that the response status code matches the expected code.
		assert.Equalf(t, test.expectedCode, postResp.StatusCode, test.description)

		// Read the response body to get the short URL.
		bodyBuffer := new(bytes.Buffer)
		bodyBuffer.ReadFrom(postResp.Body)
		shortURL := bodyBuffer.String()
		// Parse the short URL to extract the ID.
		URLParts := strings.Split(shortURL, "/")
		idURL := URLParts[len(URLParts)-1]

		// Create a GET request to retrieve the original URL.
		getReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf(`/%s`, idURL), nil)
		getReq.Header.Set("Content-Type", test.contentType)
		getResp, _ := app.Test(getReq, 1)

		// Get the redirect location header from the response.
		redirectLocation, _ := getResp.Location()
		// Assert that the redirect location matches the original URL.
		assert.Equalf(t, test.body, redirectLocation.String(), test.description)
	}
}
