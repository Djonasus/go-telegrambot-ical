package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/emersion/go-ical"
	_ "github.com/mattn/go-sqlite3"
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

func getEventsNames(fileName string) ([]ical.Event, error) {
	r, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var eve []ical.Event

	dec := ical.NewDecoder(r)
	for {
		cal, err := dec.Decode()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		for _, event := range cal.Events() {
			//summary, err := event.Props.Text(ical.PropSummary)
			if err != nil {
				return nil, err
			}
			//st, err := event.DateTimeStart(loc)
			//ed, err := event.DateTimeEnd(loc)
			eve = append(eve, event)
			//log.Printf("Found event: %v", summary)
			//log.Printf("Time start event: %v", st)
			//log.Printf("Time end event: %v", ed)
		}
	}
	return eve, nil
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

func LoadData() error {
	//tempData := []CalData{}

	db, err := sql.Open("sqlite3", "userdata.sql")
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("select * from users")
	if err != nil {
		return err
	}
	defer rows.Close()

	users = []CalData{}

	for rows.Next() {
		u := CalData{}
		id := 0
		ev := ""
		err := rows.Scan(&id, &u.userID, &u.userTime, &u.userCalendar, &u.userURL, &u.userState, &ev)
		if err != nil {
			fmt.Println(err)
			continue
		}
		_ = json.Unmarshal([]byte(ev), &u.userShowedEvents)
		users = append(users, u)
	}

	//return tempData
	return nil
}

func NewElement(data CalData) error {
	db, err := sql.Open("sqlite3", "userdata.sql")
	if err != nil {
		return err
	}
	defer db.Close()

	use_string, _ := json.Marshal(data.userShowedEvents)

	_, err = db.Exec("insert into users (userID, userTime, userCalendar, userURL, userState, userShowedEvents) values (" + strconv.FormatInt(data.userID, 10) + ", " + strconv.FormatInt(int64(data.userTime), 10) + ", '" + data.userCalendar + "', '" + data.userURL + "', '" + data.userState + "','" + string(use_string) + "')")
	if err != nil {
		return err
	}
	return nil
}

func UpdateElement(data *CalData) error {
	db, err := sql.Open("sqlite3", "userdata.sql")
	if err != nil {
		return err
	}
	defer db.Close()
	use_string, _ := json.Marshal(data.userShowedEvents)
	_, err = db.Exec("update users set userCalendar = '" + data.userCalendar + "', userTime = " + strconv.FormatInt(int64(data.userTime), 10) + ", userURL = '" + data.userURL + "', userState='" + data.userState + "', userShowedEvents = '" + string(use_string) + "'  where userID = " + strconv.FormatInt(data.userID, 10) + ";")
	if err != nil {
		return nil
	}
	return nil
}

func DeleteElement(data *CalData) error {
	db, err := sql.Open("sqlite3", "userdata.sql")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("delete from users where userID = " + strconv.FormatInt(data.userID, 10) + ";")
	if err != nil {
		return nil
	}
	return nil
}
