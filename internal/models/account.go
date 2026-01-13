package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OwnerName string         `gorm:"type:varchar(100);not null" json:"owner_name"`
	Balance   float64        `gorm:"type:decimal(15,2);not null;default:0" json:"balance"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Account) TableName() string {
	return "accounts"
}
