package repository

import (
	"sync"

	"application-design-test/internal/model"
)

type InMemoryOrderRepository struct {
	orders []model.Order
	mu     sync.Mutex
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make([]model.Order, 0),
	}
}

func (repo *InMemoryOrderRepository) SaveOrder(order *model.Order) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.orders = append(repo.orders, *order)
	return nil
}
