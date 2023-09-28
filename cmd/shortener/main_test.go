package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/ispaneli/urlpresser/internal/handlers"
	"github.com/ispaneli/urlpresser/internal/storage"
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
			route:        "/",
			body:         "https://practicum.yandex.ru/",
			contentType:  "text/plain",
			expectedCode: http.StatusCreated,
		},
		{
			description:  "Get HTTP status 400 (body was empty)",
			route:        "/",
			body:         "",
			contentType:  "text/plain",
			expectedCode: http.StatusBadRequest,
		},
		{
			description:  "Get HTTP status 201 (existing URL)",
			route:        "/",
			body:         "https://practicum.yandex.ru/",
			contentType:  "text/plain",
			expectedCode: http.StatusCreated,
		},
	}

	e := echo.New()
	h := handlers.NewHandlers(storage.NewStorage(), "http://localhost:8080/")
	e.POST("/", h.ShortingURLHandler)

	var originalURLMap = make(map[string]string)

	for _, test := range tests {
		// Create a POST request with the test data.
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(test.body))
		req.Header.Set("Content-Type", test.contentType)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Handle the request.
		if assert.NoError(t, h.ShortingURLHandler(c)) {
			// Assert that the response status code matches the expected code.
			assert.Equalf(t, test.expectedCode, rec.Code, test.description)

			// Read the response body to get the short URL.
			shortURL := rec.Body.String()

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
			route:        "/",
			body:         "https://practicum.yandex.ru/",
			contentType:  "text/plain",
			expectedCode: http.StatusCreated,
		},
		{
			description:  "Get HTTP status 201 (existing URL)",
			route:        "/",
			body:         "https://practicum.yandex.ru/",
			contentType:  "text/plain",
			expectedCode: http.StatusCreated,
		},
	}

	e := echo.New()
	h := handlers.NewHandlers(storage.NewStorage(), "http://localhost:8080/")
	e.POST("/", h.ShortingURLHandler)
	e.GET("/:id", h.RedirectingURLHandler)

	for _, test := range tests {
		// Create a POST request to shorten the URL.
		postReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(test.body))
		postReq.Header.Set("Content-Type", test.contentType)
		postResp := httptest.NewRecorder()
		e.ServeHTTP(postResp, postReq)

		// Assert that the response status code matches the expected code.
		assert.Equalf(t, test.expectedCode, postResp.Code, test.description)

		// Read the response body to get the short URL.
		shortURL := postResp.Body.String()

		// Parse the short URL to extract the ID.
		URLParts := strings.Split(shortURL, "/")
		idURL := URLParts[len(URLParts)-1]

		// Create a GET request to retrieve the original URL.
		getReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", idURL), nil)
		getReq.Header.Set("Content-Type", test.contentType)
		getResp := httptest.NewRecorder()
		e.ServeHTTP(getResp, getReq)

		// Get the redirect location header from the response.
		redirectLocation := getResp.Header().Get("Location")
		// Assert that the redirect location matches the original URL.
		assert.Equalf(t, test.body, redirectLocation, test.description)
	}
}

func TestShortenAPIHandler(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		body         string
		mainURL      string
		contentType  string
		expectedCode int
	}{
		{
			description:  "Get HTTP status 201",
			route:        "/api/shorten",
			body:         `{"url": "https://practicum.yandex.ru/"}`,
			mainURL:      "https://practicum.yandex.ru/",
			contentType:  "application/json",
			expectedCode: http.StatusCreated,
		},
		{
			description:  "Get HTTP status 201 (existing URL)",
			route:        "/api/shorten",
			body:         `{"url": "https://practicum.yandex.ru/"}`,
			mainURL:      "https://practicum.yandex.ru/",
			contentType:  "application/json",
			expectedCode: http.StatusCreated,
		},
	}

	e := echo.New()
	h := handlers.NewHandlers(storage.NewStorage(), "http://localhost:8080/")
	e.POST("/api/shorten", h.ShortenAPIHandler)
	e.GET("/:id", h.RedirectingURLHandler)

	var originalURLMap = make(map[string]string)

	for _, test := range tests {
		// Create a POST request with the test data.
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(test.body))
		req.Header.Set("Content-Type", test.contentType)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Handle the request.
		if assert.NoError(t, h.ShortenAPIHandler(c)) {
			// Assert that the response status code matches the expected code.
			assert.Equalf(t, test.expectedCode, rec.Code, test.description)

			// Read the response body to get the short URL.
			var response struct {
				Result string `json:"result"`
			}
			dataRec := rec.Body.String()
			dataRec = dataRec[1 : len(dataRec)-2]
			decodedData, err := base64.StdEncoding.DecodeString(dataRec)
			assert.NoError(t, err, test.description)
			err = json.Unmarshal(decodedData, &response)
			assert.NoError(t, err, test.description)

			// Check if the original URL already exists in the map.
			if existShortURL, ok := originalURLMap[test.body]; ok {
				// If it exists, assert that the short URL matches the existing one.
				assert.Equalf(t, existShortURL, response.Result, test.description)
			} else {
				// If it doesn't exist, add it to the map.
				originalURLMap[test.body] = response.Result
			}

			// Create a GET request to retrieve the original URL using the short URL.
			getReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", response.Result[22:]), nil)
			getReq.Header.Set("Content-Type", test.contentType)
			getResp := httptest.NewRecorder()
			e.ServeHTTP(getResp, getReq)

			// Get the redirect location header from the response.
			redirectLocation := getResp.Header().Get("Location")
			// Assert that the redirect location matches the original URL.
			assert.Equalf(t, test.mainURL, redirectLocation, test.description)
		}
	}
}
