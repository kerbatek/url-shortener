package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kerbatek/url-shortener/internal/service"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{service: service}
}

func (h *URLHandler) ShortenURL(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	url, err := h.service.Shorten(c.Request.Context(), req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, url)
}

func (h *URLHandler) RedirectURL(c *gin.Context) {
	code := c.Param("code")

	url, err := h.service.Resolve(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "short url not found"})
		return
	}

	c.Redirect(http.StatusFound, url.OriginalURL)
}

func (h *URLHandler) DeleteURL(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
