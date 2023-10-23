package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/emersion/go-ical"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// COMMAND HANDLERS
func startHandler(uid int64, bot *tgbotapi.BotAPI) {
	cd, i := FindUserById(uid, users)
	if cd.userID == 0 { //IF USER NOT EXISTS
		bot.Send(tgbotapi.NewMessage(uid, "Создание профиля..."))

		temp_data := CalData{uid, 10, "", "", "create", []string{}}
		NewElement(temp_data)
		users = append(users, temp_data)

		bot.Send(tgbotapi.NewMessage(uid, "Введите ссылку на календарь ICal:"))
	} else {
		//IF USER EXISTS
		bot.Send(tgbotapi.NewMessage(uid, "Обновление профиля..."))
		users[i].userState = "update"
		UpdateElement(&users[i])
		bot.Send(tgbotapi.NewMessage(uid, "Введите ссылку на календарь ICal:"))
	}
}

func debugHandler(uid int64, bot *tgbotapi.BotAPI) {
	cd, _ := FindUserById(uid, users)
	if cd.userID == 0 {
		bot.Send(tgbotapi.NewMessage(uid, "Профиль не обнаружен."))
		return
	}
	bot.Send(tgbotapi.NewMessage(uid, strconv.FormatInt(uid, 10)))
	bot.Send(tgbotapi.NewMessage(uid, cd.userState))
	bot.Send(tgbotapi.NewMessage(uid, cd.userCalendar))
}

func eventsHandler(uid int64, bot *tgbotapi.BotAPI) {
	cd, _ := FindUserById(uid, users)
	if cd.userID == 0 {
		bot.Send(tgbotapi.NewMessage(uid, "Профиль не обнаружен."))
		return
	}

	loc, _ := time.LoadLocation("")
	bot.Send(tgbotapi.NewMessage(uid, "События на сегодня"))

	cnt := 0

	calll, err := getEventsNames(cd.userCalendar)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, eve := range calll {
		summary, _ := eve.Props.Text(ical.PropSummary)
		curTime := time.Now()
		tm, _ := eve.DateTimeStart(loc)

		y1, m1, d1 := tm.Date()
		y2, m2, d2 := curTime.Date()

		if !(y1 == y2 && m1 == m2 && d1 == d2) {
			continue
		}
		cnt++
		mnts := fmt.Sprintf("%.0f", tm.Sub(curTime).Minutes())

		bot.Send(tgbotapi.NewMessage(uid, "Событие "+summary+". Осталось: "+(mnts)+" минут."))
	}
	if cnt == 0 {
		bot.Send(tgbotapi.NewMessage(uid, "Событий нет."))
	}
}

func syncHandler(uid int64, bot *tgbotapi.BotAPI) {
	cd, _ := FindUserById(uid, users)
	if cd.userID == 0 {
		bot.Send(tgbotapi.NewMessage(uid, "Профиль не обнаружен."))
		return
	}

	err := DownloadFile(cd.userCalendar, cd.userURL)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(uid, err.Error()))
	}

	bot.Send(tgbotapi.NewMessage(uid, "Календарь синхронизирован."))
}

func setTimeHandler(uid int64, bot *tgbotapi.BotAPI, args string) {
	if args == "" {
		bot.Send(tgbotapi.NewMessage(uid, "Укажите время, за которое бот должен предупредить вас. /settime <минуты>"))
		return
	}
	cd, i := FindUserById(uid, users)
	if cd.userID == 0 {
		bot.Send(tgbotapi.NewMessage(uid, "Профиль не обнаружен."))
		return
	}
	users[i].userTime, _ = strconv.Atoi(args)
	UpdateElement(&users[i])
	bot.Send(tgbotapi.NewMessage(uid, "Время изменено."))
}

func deleteHandler(uid int64, bot *tgbotapi.BotAPI) {
	cd, i := FindUserById(uid, users)
	if cd.userID == 0 {
		bot.Send(tgbotapi.NewMessage(uid, "Профиль не обнаружен."))
		return
	}
	DeleteElement(&cd)
	users[i] = CalData{}
	bot.Send(tgbotapi.NewMessage(uid, "Профиль удален."))
}
