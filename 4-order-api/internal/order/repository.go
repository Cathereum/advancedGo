package order

import "advancedGo/pkg/db"

type OrderRepository struct {
	OrderDB *db.DB
}

func NewOrderRepository(db *db.DB) *OrderRepository {
	return &OrderRepository{
		OrderDB: db,
	}
}

func (repo *OrderRepository) Create(order *Order) (*Order, error) {
	result := repo.OrderDB.DB.Create(order)
	if result.Error != nil {
		return nil, result.Error
	}

	return order, nil
}

func (repo *OrderRepository) GetOrderById(userId uint, id uint) (*Order, error) {
	var order Order
	result := repo.OrderDB.DB.Preload("Products").Where("user_id = ?", userId).First(&order, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &order, nil
}

func (repo *OrderRepository) GetUserOrders(userId uint) ([]Order, error) {
	var order []Order
	result := repo.OrderDB.DB.Preload("Products").Where("user_id = ?", userId).Find(&order)
	if result.Error != nil {
		return nil, result.Error
	}
	return order, nil
}
