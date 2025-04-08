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

func (p *PostgreSQL) GetOpenAuctions(ctx context.Context) ([]*models.Auction, error) {
	var auctions []*models.Auction
	err := p.db.WithContext(ctx).Where("status = ?", "open").Find(&auctions).Error
	if err != nil {
		return nil, err
	}
	return auctions, nil
}

func (p *PostgreSQL) GetAuctionByID(ctx context.Context, id string) (*models.Auction, error) {
	var auction models.Auction
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&auction).Error
	if err != nil {
		return nil, err
	}
	return &auction, nil
}
