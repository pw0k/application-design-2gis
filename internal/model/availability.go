package model

import "time"

type RoomAvailability struct {
	HotelID string
	RoomID  string
	Date    time.Time
	Quota   int
}

type DateQuotaMap map[time.Time]int

type HotelRoomKey struct {
	HotelID string
	RoomID  string
}

func NewHotelRoomKey(hotelID string, roomID string) HotelRoomKey {
	return HotelRoomKey{
		HotelID: hotelID,
		RoomID:  roomID}
}
