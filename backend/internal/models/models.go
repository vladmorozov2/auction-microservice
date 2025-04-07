package models

import (
	"time"
)

type Auction struct {
	ID              uint      `gorm:"primaryKey"`
	StartupID       int       `gorm:"not null;index"`
	Title           string    `gorm:"size:255;not null"`
	Description     *string   `gorm:"type:text"`
	StartingBid     float64   `gorm:"type:decimal(15,2);not null"`
	MinBidIncrement *float64  `gorm:"type:decimal(15,2);default:100"`
	Status          string    `gorm:"size:50;check:status IN ('open', 'closed')"`
	EndTime         time.Time `gorm:"not null"`
	WinnerID        *int      `gorm:"index"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Bid struct {
	ID         uint    `gorm:"primaryKey"`
	AuctionID  int     `gorm:"not null;index"`
	InvestorID int     `gorm:"not null;index"`
	BidAmount  float64 `gorm:"type:decimal(15,2);not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
