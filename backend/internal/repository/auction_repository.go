package repository

import (
	"context"
	"github.com/vladmorozov2/auction-service/internal/models"
)

func (p *PostgreSQL) CreateAuction(ctx context.Context, auction *models.Auction) error {
	return p.db.WithContext(ctx).Create(auction).Error
}
