package service

import (
	"application-design-test/internal/model"
	"application-design-test/internal/repository/mocks"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDaysBetween(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		from     time.Time
		to       time.Time
		expected []time.Time
	}{
		{
			name:     "same day",
			from:     time.Date(2024, 4, 27, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2024, 4, 27, 0, 0, 0, 0, time.UTC),
			expected: []time.Time{normalizeDate(time.Date(2024, 4, 27, 0, 0, 0, 0, time.UTC))},
		},
		{
			name:     "from after to",
			from:     time.Date(2024, 4, 28, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2024, 4, 27, 0, 0, 0, 0, time.UTC),
			expected: nil,
		},
		{
			name: "crossing year/month",
			from: time.Date(2023, 12, 31, 23, 0, 0, 0, time.UTC),
			to:   time.Date(2024, 1, 2, 1, 0, 0, 0, time.UTC),
			expected: []time.Time{
				normalizeDate(time.Date(2023, 12, 31, 23, 0, 0, 0, time.UTC)),
				normalizeDate(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				normalizeDate(time.Date(2024, 1, 2, 1, 0, 0, 0, time.UTC)),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := daysBetween(tt.from, tt.to)

			if tt.expected == nil {
				assert.Nil(t, result, "DaysBetween(%v, %v) должно возвращать nil", tt.from, tt.to)
			} else {
				assert.Equal(t, len(tt.expected), len(result), "DaysBetween(%v, %v) должно возвращать %v элементов, получено %v", tt.from, tt.to, len(tt.expected), len(result))
				for i, expectedDate := range tt.expected {
					assert.Equal(t, expectedDate, result[i], "День %d должен быть %v, получено %v", i, expectedDate, result[i])
				}
			}
		})
	}
}

func TestNormalizeDate(t *testing.T) {
	t.Parallel()
	input := time.Date(2024, 4, 27, 15, 30, 45, 123456789, time.UTC)
	expected := time.Date(2024, 4, 27, 0, 0, 0, 0, time.UTC)

	result := normalizeDate(input)

	assert.Equal(t, expected, result, "NormalizeDate(%v) должно возвращать %v, получено %v", input, expected, result)
}

func TestBookingService_CreateOrder(t *testing.T) {
	t.Parallel()
	order := &model.Order{
		HotelID: "reddison",
		RoomID:  "lux",
		Email:   "guest@mail.ru",
		From:    time.Date(2023, 12, 30, 0, 0, 0, 0, time.UTC),
		To:      time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
	}
	bookingInterval := daysBetween(order.From, order.To)
	initialQuota := model.DateQuotaMap{
		time.Date(2023, 12, 30, 0, 0, 0, 0, time.UTC): 2,
		time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC): 2,
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC):   1,
		time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC):   1,
	}

	ctx := context.Background()
	mockOrderRepo := mocks.NewOrderRepositoryMock(t)
	mockAvailabilityRepo := mocks.NewAvailabilityRepositoryMock(t)
	mockAvailabilityRepo.GetRoomAvailabilityMock.
		When(ctx, order.HotelID, order.RoomID).
		Then(initialQuota, nil)
	mockOrderRepo.SaveOrderMock.
		When(ctx, order).
		Then(nil)
	mockAvailabilityRepo.DecrementRoomQuotaMock.
		When(ctx, order.HotelID, order.RoomID, bookingInterval).
		Then(nil)
	service := NewBookingService(mockOrderRepo, mockAvailabilityRepo)

	err := service.CreateOrder(ctx, order)

	assert.NoError(t, err, "CreateOrder shouldn't return error")
	assert.True(t, mockAvailabilityRepo.MinimockGetRoomAvailabilityDone(), "GetRoomAvailability должен быть вызван один раз")
	assert.True(t, mockOrderRepo.MinimockSaveOrderDone(), "SaveOrder должен быть вызван один раз")
	assert.True(t, mockAvailabilityRepo.MinimockDecrementRoomQuotaDone(), "DecrementRoomQuota должен быть вызван один раз")
}

func TestBookingService_CreateOrder_QuotaError(t *testing.T) {
	t.Parallel()
	order := &model.Order{
		HotelID: "reddison",
		RoomID:  "lux",
		Email:   "guest@mail.ru",
		From:    time.Date(2023, 12, 30, 0, 0, 0, 0, time.UTC),
		To:      time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
	}
	// квота 0
	insufficientQuota := model.DateQuotaMap{
		time.Date(2023, 12, 30, 0, 0, 0, 0, time.UTC): 0,
		time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC): 0,
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC):   0,
	}

	ctx := context.Background()
	mockOrderRepo := mocks.NewOrderRepositoryMock(t)
	mockAvailabilityRepo := mocks.NewAvailabilityRepositoryMock(t)
	mockAvailabilityRepo.
		GetRoomAvailabilityMock.
		When(ctx, order.HotelID, order.RoomID).
		Then(insufficientQuota, nil)
	service := NewBookingService(mockOrderRepo, mockAvailabilityRepo)

	err := service.CreateOrder(ctx, order)

	assert.Error(t, err, "CreateOrder should return error")
	assert.Contains(t, err.Error(), "qota is unavailable for specific date")
	assert.True(t, mockAvailabilityRepo.MinimockGetRoomAvailabilityDone(), "GetRoomAvailability должен быть вызван один раз")
	assert.Equal(t, uint64(0), mockOrderRepo.SaveOrderAfterCounter(), "SaveOrder не должен вызываться при недостаточной квоте")
	assert.Equal(t, uint64(0), mockAvailabilityRepo.DecrementRoomQuotaAfterCounter(), "DecrementRoomQuota не должен вызываться при недостаточной квоте")
}
