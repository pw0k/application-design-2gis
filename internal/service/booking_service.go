package service

import (
	"context"
	"fmt"
	"time"

	"application-design-test/internal/model"
)

type AvailabilityRepository interface {
	GetRoomAvailability(ctx context.Context, hotelID, roomID string) (model.DateQuotaMap, error)
	DecrementRoomQuota(ctx context.Context, hotelID, roomID string, date []time.Time) error
}

type OrderRepository interface {
	SaveOrder(ctx context.Context, order *model.Order) error
}

type BookingService struct {
	orderRepo        OrderRepository
	availabilityRepo AvailabilityRepository
}

func NewBookingService(orderRepo OrderRepository, availabilityRepo AvailabilityRepository) *BookingService {
	return &BookingService{
		orderRepo:        orderRepo,
		availabilityRepo: availabilityRepo,
	}
}

func (s *BookingService) CreateOrder(ctx context.Context, order *model.Order) error {
	dateQuotaMap, err := s.availabilityRepo.GetRoomAvailability(ctx, order.HotelID, order.RoomID)
	if err != nil {
		return fmt.Errorf("GetRoomAvailability error %v, %w", order, err)
	}
	for currDate := order.From; !currDate.After(order.To); currDate = currDate.AddDate(0, 0, 1) {
		currQuote, ok := dateQuotaMap[currDate]
		if !ok || currQuote < 1 {
			return newQuotaUnavailableError(order)
		}
	}

	//todo: здесь по хорошему нужна транзакция
	if err := s.orderRepo.SaveOrder(ctx, order); err != nil {
		return fmt.Errorf("SaveOrder error, order %v, %w", order, err)
	}
	bookingInterval := daysBetween(order.From, order.To)
	if err := s.availabilityRepo.DecrementRoomQuota(ctx, order.HotelID, order.RoomID, bookingInterval); err != nil {
		return fmt.Errorf("DecrementRoomQuota error, order %v, %w", order, err)
	}

	return nil
}

func daysBetween(from time.Time, to time.Time) []time.Time {
	if from.After(to) {
		return nil
	}

	days := make([]time.Time, 0)
	for d := normalizeDate(from); !d.After(normalizeDate(to)); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}
	return days
}

func normalizeDate(timestamp time.Time) time.Time {
	return time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.UTC)
}
