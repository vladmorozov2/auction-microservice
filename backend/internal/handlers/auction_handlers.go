package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vladmorozov2/auction-service/internal/models"
	"github.com/vladmorozov2/auction-service/internal/repository"
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
	fmt.Println("Creating auction in handlers:", auction)

	if auction.StartingBid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "starting_bid must be greater than 0"})
		return
	}
	if auction.MinBidIncrement != nil && *auction.MinBidIncrement <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_bid_increment must be greater than 0"})
		return
	}
	if auction.EndTime.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_time is required"})
		return
	}
	if auction.EndTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_time must be in the future"})
		return
	}

	if err := h.repo.CreateAuction(c.Request.Context(), &auction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create auction"})
		return
	}

	c.JSON(http.StatusCreated, auction)
}
