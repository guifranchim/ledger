package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionTypeDebit   TransactionType = "DEBIT"
	TransactionTypeCredit  TransactionType = "CREDIT"
	TransactionTypeReverse TransactionType = "REVERSE"
)

type Transaction struct {
	ID          string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountID   string          `gorm:"type:uuid;not null;index" json:"account_id"`
	Type        TransactionType `gorm:"type:varchar(20);not null" json:"type"`
	Amount      float64         `gorm:"type:decimal(15,2);not null" json:"amount"`
	Description string          `gorm:"type:varchar(255)" json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`

	Account Account `gorm:"foreignKey:AccountID" json:"-"`
}

func (Transaction) TableName() string {
	return "transactions"
}
