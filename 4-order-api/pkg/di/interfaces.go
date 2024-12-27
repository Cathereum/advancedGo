package di

import (
	"advancedGo/internal/product"
)

type IProductRepository interface {
	GetProductById(id uint) (*product.Product, error)
}

type IUserRepository interface {
	GetUserId(phone string) (userId uint, err error)
}
