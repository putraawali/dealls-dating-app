package models

import (
	"time"
)

type Swipe struct {
	SwipeID      int64     `json:"swipe_id" gorm:"primaryKey;column:swipe_id"`
	UserID       int64     `json:"user_id" gorm:"column:user_id"`
	TargetUserID int64     `json:"target_user_id" gorm:"column:target_user_id"`
	Pass         bool      `json:"pass" gorm:"column:pass"`
	Like         bool      `json:"like" gorm:"column:like"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`

	// Association
	User   User `json:"user" gorm:"association_foreignkey:user_id"`
	Target User `json:"target" gorm:"association_foreignkey:target_user_id"`
}
