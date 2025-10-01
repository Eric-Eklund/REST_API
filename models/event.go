package models

import "time"

type Event struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	DateTime    time.Time `json:"date_time" binding:"required"`
	UserID      int       `json:"user_id"`
}

var events []Event

func (e *Event) Save() error {
	events = append(events, *e)
	return nil
}

func GetAllEvents() []Event {
	return events
}
