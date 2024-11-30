package model

import "time"

type Order struct {
	HotelID string    `json:"hotel_id"`
	RoomID  string    `json:"room_id"`
	Email   string    `json:"email"`
	From    time.Time `json:"from"`
	To      time.Time `json:"to"`
}
