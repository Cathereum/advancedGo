package order

import (
	"advancedGo/internal/product"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserId   uint              `json:"user_id"`
	Products []product.Product `gorm:"many2many:order_products;"`
}
