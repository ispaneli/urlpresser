package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ispaneli/urlpresser/internal/storage"
)

type Handlers struct {
	baseURL string
	storage *storage.Store
}

type shortenRequest struct {
	URL string `json:"url"`
}

func (h *Handlers) ShortingURLHandler(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	shortURL := h.storage.GetShortURL(string(body))
	return c.String(http.StatusCreated, fmt.Sprintf(`%s%s`, h.baseURL, shortURL))
}

func (h *Handlers) RedirectingURLHandler(c echo.Context) error {
	id := c.Param("id")

	originalURL, ok := h.storage.GetOriginURL(id)
	if !ok {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	c.Response().Header().Set("Location", originalURL)
	return c.NoContent(http.StatusTemporaryRedirect)
}

func (h *Handlers) ShortenAPIHandler(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	var req shortenRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		return err
	}

	shortURL := h.storage.GetShortURL(req.URL)

	resJSON, err := json.Marshal(map[string]string{"result": fmt.Sprintf(`%s%s`, h.baseURL, shortURL)})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, resJSON)
}

func NewHandlers(storage *storage.Store, baseURL string) Handlers {
	return Handlers{
		baseURL: baseURL,
		storage: storage,
	}
}
