package repository

import (
	"github.com/vladmorozov2/auction-service/internal/models"
	"gorm.io/gorm"
	// "time"
	"context"
)

type Repository interface {
	// GetAuctionByID(id int) (*models.Auction, error)
	CreateAuction(ctx context.Context, auction *models.Auction) error
}

type PostgreSQL struct {
	db *gorm.DB
}

func NewAuctionRepository(db *gorm.DB) Repository {
	return &PostgreSQL{db: db}
}
