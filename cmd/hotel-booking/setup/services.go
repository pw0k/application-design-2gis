package setup

import "application-design-test/internal/service"

type Services struct {
	BookingService service.BookingService
}

func InitServices(repos *Repositories) *Services {
	bookingService := service.NewBookingService(repos.OrderRepo, repos.AvailabilityRepo)
	return &Services{
		BookingService: bookingService,
	}
}
