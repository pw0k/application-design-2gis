package repository

import (
	"application-design-test/internal/model"
	"sync"
)

type OrderRepository interface {
	SaveOrder(order *model.Order) error
}

type inMemoryOrderRepository struct {
	orders []model.Order
	mu     sync.Mutex
}

func NewInMemoryOrderRepository() OrderRepository {
	return &inMemoryOrderRepository{
		orders: make([]model.Order, 0),
	}
}

func (repo *inMemoryOrderRepository) SaveOrder(order *model.Order) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.orders = append(repo.orders, *order)
	return nil
}
