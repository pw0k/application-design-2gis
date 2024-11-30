package service

import (
	"application-design-test/internal/model"
	"application-design-test/internal/repository"
	"fmt"
	"time"
)

type BookingService interface {
	CreateOrder(order *model.Order) error
}

type bookingService struct {
	orderRepo        repository.OrderRepository
	availabilityRepo repository.AvailabilityRepository
}

func NewBookingService(orderRepo repository.OrderRepository, availabilityRepo repository.AvailabilityRepository) BookingService {
	return &bookingService{
		orderRepo:        orderRepo,
		availabilityRepo: availabilityRepo,
	}
}

func (s *bookingService) CreateOrder(order *model.Order) error {
	dateQuotaMap, err := s.availabilityRepo.GetRoomAvailability(order.HotelID, order.RoomID)
	if err != nil {
		return fmt.Errorf("проблемы при получении бронирований %v, %w", order, err)
	}
	for currDate := order.From; currDate.Before(order.To); currDate = currDate.AddDate(0, 0, 1) {
		currQuote, ok := dateQuotaMap[currDate]
		if !ok || currQuote < 1 {
			return fmt.Errorf("номер не доступен в выбранные даты, order %v", order)
		}
	}

	if err := s.orderRepo.SaveOrder(order); err != nil {
		return fmt.Errorf("проблемы при сохранении, order %v, %w", order, err)
	}
	bookingInterval := DaysBetween(order.From, order.To)
	if err := s.availabilityRepo.DecrementRoomQuota(order.HotelID, order.RoomID, bookingInterval); err != nil {
		return fmt.Errorf("не получилось обновить бронирования, order %v, %w", order, err)
	}

	return nil
}

func DaysBetween(from time.Time, to time.Time) []time.Time {
	if from.After(to) {
		return nil
	}

	days := make([]time.Time, 0)
	for d := NormalizeDate(from); !d.After(NormalizeDate(to)); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}
	return days
}

func NormalizeDate(timestamp time.Time) time.Time {
	return time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.UTC)
}
