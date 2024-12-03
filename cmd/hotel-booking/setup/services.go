package setup

import (
	"application-design-test/internal/handler"
	"application-design-test/internal/service"
)

type Services struct {
	BookingService handler.BookingService
}

func InitServices(repos *Repositories) *Services {
	bookingService := service.NewBookingService(repos.OrderRepo, repos.AvailabilityRepo)
	return &Services{
		BookingService: bookingService,
	}
}
