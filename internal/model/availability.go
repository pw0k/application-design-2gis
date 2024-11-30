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
