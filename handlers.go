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
		users = append(users, CalData{uid, "", "", "create", []string{}})
		bot.Send(tgbotapi.NewMessage(uid, "Введите ссылку на календарь ICal:"))
	} else {
		//IF USER EXISTS
		bot.Send(tgbotapi.NewMessage(uid, "Обновление профиля..."))
		users[i].userState = "update"
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

func eventsHandler(uid int64, bot *tgbotapi.BotAPI, args string) {
	fmt.Println(args)
	cd, _ := FindUserById(uid, users)
	if cd.userID == 0 {
		bot.Send(tgbotapi.NewMessage(uid, "Профиль не обнаружен."))
		return
	}

	loc, _ := time.LoadLocation("")

	for _, eve := range getEventsNames(cd.userCalendar) {
		summary, _ := eve.Props.Text(ical.PropSummary)
		curTime := time.Now()
		tm, _ := eve.DateTimeStart(loc)

		mnts := fmt.Sprintf("%.0f", tm.Sub(curTime).Minutes())

		bot.Send(tgbotapi.NewMessage(uid, "Событие "+summary+". Осталось: "+(mnts)+" минут."))
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
