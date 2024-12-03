package setup

import (
	"application-design-test/internal/repository"
	"application-design-test/internal/service"
)

type Repositories struct {
	OrderRepo        service.OrderRepository
	AvailabilityRepo service.AvailabilityRepository
}

func InitRepositories() *Repositories {
	orderRepo := repository.NewInMemoryOrderRepository()
	availabilityRepo := repository.NewInMemoryAvailabilityRepository()

	return &Repositories{
		OrderRepo:        orderRepo,
		AvailabilityRepo: availabilityRepo,
	}
}
