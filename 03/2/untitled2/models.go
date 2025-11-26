
// models.go
package main

import (
	"time"
)

type Account struct {
	ID      uint    `gorm:"primaryKey"`
	Balance float64 `gorm:"type:decimal(10,2);default:0.00"`
}

type Transaction struct {
	ID            uint           `gorm:"primaryKey"`
	FromAccountID uint           `gorm:"not null"`
	ToAccountID   uint           `gorm:"not null"`
	Amount        float64        `gorm:"type:decimal(10,2);not null"`
	CreatedAt     time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
}