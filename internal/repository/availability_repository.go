package repository

import (
	"application-design-test/internal/model"
	"fmt"
	"sync"
	"time"
)

type AvailabilityRepository interface {
	GetRoomAvailability(hotelID, roomID string) (model.DateQuotaMap, error)
	DecrementRoomQuota(hotelID, roomID string, date []time.Time) error
}

type inMemoryAvailabilityRepository struct {
	availabilityMap map[model.HotelRoomKey]model.DateQuotaMap
	mu              sync.Mutex
}

func NewInMemoryAvailabilityRepository() AvailabilityRepository {
	InitialAvailability := []model.RoomAvailability{
		{"reddison", "lux", newDate(2024, 1, 1), 1},
		{"reddison", "lux", newDate(2024, 1, 2), 1},
		{"reddison", "lux", newDate(2024, 1, 3), 1},
		{"reddison", "lux", newDate(2024, 1, 4), 1},
		{"reddison", "lux", newDate(2024, 1, 5), 0},
	}

	roomAvailabilityMap := make(map[model.HotelRoomKey]model.DateQuotaMap)
	for _, avail := range InitialAvailability {
		key := model.HotelRoomKey{HotelID: avail.HotelID, RoomID: avail.RoomID}
		_, ok := roomAvailabilityMap[key]
		if !ok {
			roomAvailabilityMap[key] = model.DateQuotaMap{}
		}
		roomAvailabilityMap[key][avail.Date] = avail.Quota
	}

	return &inMemoryAvailabilityRepository{
		availabilityMap: roomAvailabilityMap,
	}
}

func (repo *inMemoryAvailabilityRepository) GetRoomAvailability(hotelID, roomID string) (model.DateQuotaMap, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	hotelRoomKey := model.HotelRoomKey{HotelID: hotelID, RoomID: roomID}
	dateQuotaMap, ok := repo.availabilityMap[hotelRoomKey]
	if !ok {
		return dateQuotaMap, fmt.Errorf("бронь для номера не доступна, hotelID %v, roomID %v", hotelID, roomID)
	}
	//todo: передавать копию ?
	return dateQuotaMap, nil
}

func (repo *inMemoryAvailabilityRepository) DecrementRoomQuota(hotelID, roomID string, bookingDate []time.Time) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	hotelRoomKey := model.HotelRoomKey{HotelID: hotelID, RoomID: roomID}
	dateQuotaMap, ok := repo.availabilityMap[hotelRoomKey]
	if !ok {
		return fmt.Errorf("бронь для номера не доступна, hotelID %v, roomID %v", hotelID, roomID)
	}
	for _, currDate := range bookingDate {
		quota, ok := dateQuotaMap[currDate]
		if !ok || quota < 1 {
			return fmt.Errorf("бронь на конкретную дату для номера не доступна, hotelID %v, roomID %v, currDate %v, quota %v", hotelID, roomID, currDate, quota)
		}
		dateQuotaMap[currDate] = quota - 1
	}
	return nil
}

func newDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
