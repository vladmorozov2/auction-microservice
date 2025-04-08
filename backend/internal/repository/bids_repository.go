package repository

import (
	"context"

	"fmt"
	"github.com/google/uuid"
	"github.com/vladmorozov2/auction-service/internal/models"

	"time"
)

func (p *PostgreSQL) CreateBid(ctx context.Context, bid *models.Bid) error {
	bid.UpdatedAt = time.Now()
	bid.CreatedAt = time.Now()
	bid.ID = uuid.New()

	fmt.Println("Creating bid:", bid)

	return p.db.WithContext(ctx).Create(bid).Error
}

func (p *PostgreSQL) GetLastBidForAuction(ctx context.Context, auctionID string) (*models.Bid, error) {
	var bid models.Bid
	err := p.db.WithContext(ctx).Where("auction_id = ?", auctionID).Order("created_at desc").First(&bid).Error
	if err != nil {
		return nil, err
	}
	return &bid, nil
}
