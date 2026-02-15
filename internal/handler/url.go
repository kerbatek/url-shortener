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
	// TODO: parse request body, call service.Shorten, return JSON response
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (h *URLHandler) RedirectURL(c *gin.Context) {
	// TODO: extract code param, call service.Resolve, redirect
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (h *URLHandler) DeleteURL(c *gin.Context) {
	// TODO: extract ID param, call service.Delete, return status
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
