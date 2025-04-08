package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vladmorozov2/auction-service/internal/models"
	"github.com/vladmorozov2/auction-service/internal/repository"
	"net/http"
)

func (h *Handler) CreateBid(c *gin.Context) {
	auctionIDStr := c.Param("auction_id")
	if auctionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "auction_id is required in URL"})
		return
	}

	auctionID, err := uuid.Parse(auctionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid auction_id format"})
		return
	}

	var bid models.Bid
	if err := c.ShouldBindJSON(&bid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	bid.AuctionID = auctionID

	if bid.BidAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bid_amount must be greater than 0"})
		return
	}
	if bid.InvestorID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "investor_id must be greater than 0"})
		return
	}

	auction, err := h.repo.GetAuctionByID(c.Request.Context(), auctionID.String())
	if err != nil {
		if errors.Is(err, repository.ErrAuctionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get auction"})
		return
	}
	if auction.Status != "open" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "auction is not open"})
		return
	}

	if bid.BidAmount < auction.StartingBid {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("bid must be greater than or equal to the starting bid of %d", auction.StartingBid)})
		return
	}

	lastBid, err := h.repo.GetLastBidForAuction(c.Request.Context(), auctionID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get last bid"})
		return
	}

	if lastBid != nil {
		minIncrement := *auction.MinBidIncrement
		if bid.BidAmount < lastBid.BidAmount+minIncrement {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("bid must be greater than the last bid plus the minimum increment of %d", minIncrement)})
			return
		}
	}
	if err := h.repo.CreateBid(c.Request.Context(), &bid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create bid"})
		return
	}

	c.JSON(http.StatusCreated, bid)
}
