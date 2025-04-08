package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/vladmorozov2/auction-service/internal/models"
	"time"
)

func (p *PostgreSQL) CreateAuction(ctx context.Context, auction *models.Auction) error {
	auction.UpdatedAt = time.Now()
	auction.CreatedAt = time.Now()
	auction.ID = uuid.New()
	if auction.Status == "" {
		auction.Status = "open"
	}

	fmt.Println("Creating auction:", auction)

	return p.db.WithContext(ctx).Create(auction).Error
}
