package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	TransactionID int64     `json:"transaction_id" gorm:"primaryKey;column:transaction_id"`
	UserID        int64     `json:"user_id" gorm:"column:user_id"`
	Status        string    `json:"status" gorm:"column:status;type:varchar(100)"`
	PaymentMethod string    `json:"payment_method" gorm:"column:payment_method;type:varchar(100)"`
	VaNumber      string    `json:"va_number" gorm:"column:va_number;type:varchar(50)"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at"`

	User User `json:"user" gorm:"association_foreignkey:user_id"`
}

func (t *Transaction) BeforeCreate(*gorm.DB) (err error) {
	t.Status = "pending"
	return
}
