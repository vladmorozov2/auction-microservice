package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/vladmorozov2/auction-service/internal/models"
	"github.com/vladmorozov2/auction-service/internal/repository"
	"net/http"
)

type Handler struct {
	repo repository.Repository
}

func NewHandler(repo repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) CreateAuction(c *gin.Context) {
	var auction models.Auction

	if err := c.ShouldBindJSON(&auction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateAuction(c.Request.Context(), &auction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create auction"})
		return
	}

	c.JSON(http.StatusCreated, auction)
}
