package models

import "gorm.io/gorm"

// Ticket model for database
type Ticket struct {
	gorm.Model        // Embeds ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string `json:"name" gorm:"not null"`
	Email      string `json:"email" gorm:"not null;uniqueIndex"`
	Tickets    uint   `json:"tickets" gorm:"not null;default:1"`
}
