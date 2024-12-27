package user

import (
	"advancedGo/internal/order"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Phone            string        `gorm:"index"`
	SessionId        string        `json:"session_id"`
	VerificationCode string        `json:"verification_code"`
	IsVerified       bool          `json:"is_verified"`
	Orders           []order.Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
