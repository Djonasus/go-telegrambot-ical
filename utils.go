package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/emersion/go-ical"
)

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func ExampleDecoder(fileName string) {
	// Let's assume r is an io.Reader containing iCal data
	//var r io.Reader

	r, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	loc, _ := time.LoadLocation("")

	dec := ical.NewDecoder(r)
	for {
		cal, err := dec.Decode()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		for _, event := range cal.Events() {
			summary, err := event.Props.Text(ical.PropSummary)
			if err != nil {
				log.Fatal(err)
			}
			st, err := event.DateTimeStart(loc)
			ed, err := event.DateTimeEnd(loc)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Found event: %v", summary)
			log.Printf("Time start event: %v", st)
			log.Printf("Time end event: %v", ed)
		}
	}
}

func getEventsNames(fileName string) []ical.Event {
	r, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	var eve []ical.Event

	dec := ical.NewDecoder(r)
	for {
		cal, err := dec.Decode()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		for _, event := range cal.Events() {
			//summary, err := event.Props.Text(ical.PropSummary)
			if err != nil {
				log.Fatal(err)
			}
			//st, err := event.DateTimeStart(loc)
			//ed, err := event.DateTimeEnd(loc)
			eve = append(eve, event)
			//log.Printf("Found event: %v", summary)
			//log.Printf("Time start event: %v", st)
			//log.Printf("Time end event: %v", ed)
		}
	}
	return eve
}

func FindUserById(uid int64, utables []CalData) (CalData, int) {
	if len(utables) == 0 {
		return CalData{}, 0
	}
	for i, cd := range utables {
		if cd.userID == uid {
			return cd, i
		}
	}
	return CalData{}, 0
}
