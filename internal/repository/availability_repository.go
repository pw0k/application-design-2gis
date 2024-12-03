package repository

import (
	"fmt"
	"sync"
	"time"

	"application-design-test/internal/model"
)

type InMemoryAvailabilityRepository struct {
	availabilityMap map[model.HotelRoomKey]model.DateQuotaMap
	mu              sync.RWMutex
}

func NewInMemoryAvailabilityRepository() *InMemoryAvailabilityRepository {
	initialAvailability := []model.RoomAvailability{
		{"reddison", "lux", newDate(2024, 1, 1), 1},
		{"reddison", "lux", newDate(2024, 1, 2), 1},
		{"reddison", "lux", newDate(2024, 1, 3), 1},
		{"reddison", "lux", newDate(2024, 1, 4), 1},
		{"reddison", "lux", newDate(2024, 1, 5), 0},
	}

	roomAvailabilityMap := make(map[model.HotelRoomKey]model.DateQuotaMap)
	for _, avail := range initialAvailability {
		key := model.NewHotelRoomKey(avail.HotelID, avail.RoomID)
		_, ok := roomAvailabilityMap[key]
		if !ok {
			roomAvailabilityMap[key] = model.DateQuotaMap{}
		}
		roomAvailabilityMap[key][avail.Date] = avail.Quota
	}

	return &InMemoryAvailabilityRepository{
		availabilityMap: roomAvailabilityMap,
	}
}

func (repo *InMemoryAvailabilityRepository) GetRoomAvailability(hotelID, roomID string) (model.DateQuotaMap, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	hotelRoomKey := model.NewHotelRoomKey(hotelID, roomID)
	dateQuotaMap, ok := repo.availabilityMap[hotelRoomKey]
	if !ok {
		return dateQuotaMap, fmt.Errorf("dateQuotaMap is unavailable, hotelID %v, roomID %v", hotelID, roomID)
	}
	return dateQuotaMap, nil
}

func (repo *InMemoryAvailabilityRepository) DecrementRoomQuota(hotelID, roomID string, bookingDate []time.Time) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	hotelRoomKey := model.NewHotelRoomKey(hotelID, roomID)
	dateQuotaMap, ok := repo.availabilityMap[hotelRoomKey]
	if !ok {
		return fmt.Errorf("dateQuotaMap is unavailable, hotelID %v, roomID %v", hotelID, roomID)
	}
	for _, currDate := range bookingDate {
		quota, ok := dateQuotaMap[currDate]
		if !ok || quota < 1 {
			return fmt.Errorf("qota is unavailable for specific date, hotelID %v, roomID %v, currDate %v, quota %v", hotelID, roomID, currDate, quota)
		}
		dateQuotaMap[currDate] = quota - 1
	}
	return nil
}

func newDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
