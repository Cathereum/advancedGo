package product

import (
	"advancedGo/pkg/db"

	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	ProductDB *db.DB
}

func NewProductRepository(db *db.DB) *ProductRepository {
	return &ProductRepository{
		ProductDB: db,
	}
}

func (repo *ProductRepository) Create(product *Product) (*Product, error) {
	result := repo.ProductDB.DB.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (repo *ProductRepository) GetProductById(id uint) (*Product, error) {
	var product Product
	result := repo.ProductDB.DB.First(&product, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &product, nil
}

func (repo *ProductRepository) DeleteProductById(id uint) error {
	result := repo.ProductDB.DB.Delete(&Product{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ProductRepository) UpdateProductById(product *Product) (*Product, error) {
	result := repo.ProductDB.DB.Clauses(clause.Returning{}).Updates(product)
	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (repo *ProductRepository) GetAllProducts() ([]Product, error) {
	var products []Product
	result := repo.ProductDB.DB.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}
