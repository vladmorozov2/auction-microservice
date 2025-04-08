package repository

import (
	"context"
	"github.com/vladmorozov2/auction-service/internal/models"
	"fmt"
)

func (p *PostgreSQL) CreateAuction(ctx context.Context, auction *models.Auction) error {
	fmt.Println("Creating auction:", auction)
	return p.db.WithContext(ctx).Create(auction).Error
}
