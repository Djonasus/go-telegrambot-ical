package main

import "time"

type CalData struct {
	userID   int64
	userTime int
	//userCalendar string
	userURL    string
	userState  string
	userEvents []EventData
	//userShowedEvents []string
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
