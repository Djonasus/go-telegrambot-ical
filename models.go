package main

import "time"

type CalData struct {
	userID     int64
	userTime   int
	userURL    string
	userState  string
	userEvents []EventData
}

type JsonElement struct {
	Name     string
	Class    string
	Text     string
	Elements []JsonElement
}

type EventData struct {
	Name   string
	Date   time.Time
	Showed bool
}
