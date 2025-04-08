package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vladmorozov2/auction-service/internal/models"
	"github.com/vladmorozov2/auction-service/internal/repository"
	"gorm.io/gorm"
	"net/http"
	"time"
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

func (h *Handler) GetOpenAuctions(c *gin.Context) {
	auctions, err := h.repo.GetOpenAuctions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get open auctions"})
		return
	}

	c.JSON(http.StatusOK, auctions)
}
func (h *Handler) GetAuctionByID(c *gin.Context) {
	id := c.Param("id")

	// Optional: Validate UUID format (if using a UUID library)
	// For simplicity, we'll assume it's a valid UUID string here
	auction, err := h.repo.GetAuctionByID(c.Request.Context(), id)
	if err != nil {
		// Check if the error is "not found"
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "auction not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve auction"})
		}
		return
	}

	c.JSON(http.StatusOK, auction)
}
