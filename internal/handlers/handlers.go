package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ispaneli/urlpresser/internal/storage"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	baseURL string
	storage *storage.Store
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

func NewHandlers(storage *storage.Store, baseURL string) Handlers {
	return Handlers{
		baseURL: baseURL,
		storage: storage,
	}
}
