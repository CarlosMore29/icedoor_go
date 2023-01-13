package model

import "time"

type User struct {
	ID   int
	Name string
	City string
}

type Timeline struct {
	ID   string    `json:"id"`
	Data string    `json:"data"`
	Date time.Time `json:"date"`
}
