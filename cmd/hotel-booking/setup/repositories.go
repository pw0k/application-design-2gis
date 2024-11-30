package setup

import (
	"application-design-test/internal/repository"
)

type Repositories struct {
	OrderRepo        repository.OrderRepository
	AvailabilityRepo repository.AvailabilityRepository
}

func InitRepositories() *Repositories {
	orderRepo := repository.NewInMemoryOrderRepository()
	availabilityRepo := repository.NewInMemoryAvailabilityRepository()

	return &Repositories{
		OrderRepo:        orderRepo,
		AvailabilityRepo: availabilityRepo,
	}
}
