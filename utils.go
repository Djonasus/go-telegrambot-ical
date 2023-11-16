package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/goodsign/monday"
	"github.com/itzg/restify"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/html"
)

func GetEvents(cd CalData) ([]EventData, error) {

	tempEvents := []EventData{}

	uri, _ := url.Parse(cd.userURL)
	q := uri.Query()
	q.Set("limit", MAX_EVENTS)
	uri.RawQuery = q.Encode()
	ht, _ := restify.LoadContent(uri, "")
	dat, _ := restify.ConvertHtmlToJson([]*html.Node{ht.LastChild.LastChild.LastChild})

	var arr []JsonElement

	err := json.Unmarshal(dat, &arr)
	if err != nil {
		return nil, err
	}
	for _, v := range arr[0].Elements {
		newEvent := EventData{}
		newEvent.Name = v.Elements[0].Text

		newstr, err := strconv.Unquote(`"` + v.Elements[1].Elements[0].Text + `"`)
		if err != nil {
			return nil, err
		}
		loc, _ := time.LoadLocation("Local")
		tmg, err := monday.ParseInLocation("2 January 2006 15:04", newstr, loc, monday.LocaleRuRU)
		if err != nil {
			return nil, err
		}

		newEvent.Date = tmg

		tempEvents = append(tempEvents, newEvent)
	}

	return tempEvents, nil
}

func LoadData() error {
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

	users = make(map[int64]*CalData)

	for rows.Next() {
		u := CalData{}

		ev := ""
		err := rows.Scan(&u.userID, &u.userTime, &u.userURL, &u.userState, &ev)
		if err != nil {
			fmt.Println(err)
			continue
		}
		_ = json.Unmarshal([]byte(ev), &u.userEvents)
		users[u.userID] = &u
	}
	return nil
}

func NewElement(data CalData) error {
	db, err := sql.Open("sqlite3", "userdata.sql")
	if err != nil {
		return err
	}
	defer db.Close()

	use_string, _ := json.Marshal(data.userEvents)

	_, err = db.Exec("insert into users (userID, userTime, userURL, userState, userEvents) values (" + strconv.FormatInt(data.userID, 10) + ", " + strconv.FormatInt(int64(data.userTime), 10) + ", '" + data.userURL + "', '" + data.userState + "','" + string(use_string) + "')")
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
	use_string, err := json.Marshal(data.userEvents)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = db.Exec("update users set userTime = " + strconv.FormatInt(int64(data.userTime), 10) + ", userURL = '" + data.userURL + "', userState='" + data.userState + "', userEvents = '" + string(use_string) + "'  where userID = " + strconv.FormatInt(data.userID, 10) + ";")
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

func removeEvent(slice []EventData, s int) []EventData {
	return append(slice[:s], slice[s+1:]...)
}
