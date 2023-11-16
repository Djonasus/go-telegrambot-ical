package main

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// COMMAND HANDLERS
func startHandler(uid int64, bot *tgbotapi.BotAPI) {
	if cd, ok := users[uid]; ok {
		//IF USER EXISTS
		bot.Send(tgbotapi.NewMessage(uid, "Обновление профиля..."))
		cd.userState = "update"
		UpdateElement(cd)
		bot.Send(tgbotapi.NewMessage(uid, "Введите ссылку на календарь ICal:"))
	} else {
		bot.Send(tgbotapi.NewMessage(uid, "Создание профиля..."))
		temp_data := CalData{uid, 10, "", "create", nil}
		NewElement(temp_data)
		users[uid] = &temp_data

		bot.Send(tgbotapi.NewMessage(uid, "Введите ссылку на календарь HTML:"))
	}
}

func debugHandler(uid int64, bot *tgbotapi.BotAPI) {
	if cd, ok := users[uid]; ok {
		bot.Send(tgbotapi.NewMessage(uid, strconv.FormatInt(uid, 10)))
		bot.Send(tgbotapi.NewMessage(uid, cd.userState))
	} else {
		bot.Send(tgbotapi.NewMessage(uid, "Профиль не обнаружен."))
	}
}

func eventsHandler(uid int64, bot *tgbotapi.BotAPI) {
	if cd, ok := users[uid]; ok {
		bot.Send(tgbotapi.NewMessage(uid, "Ближайшие 5 событий:"))

		cnt := 0
		for _, eve := range cd.userEvents {
			bot.Send(tgbotapi.NewMessage(uid, "Событие: "+eve.Name+". Время: "+eve.Date.String()))
			cnt++
		}
		if cnt == 0 {
			bot.Send(tgbotapi.NewMessage(uid, "Событий нет."))
		}
	} else {
		bot.Send(tgbotapi.NewMessage(uid, "Профиль не обнаружен."))
	}
}

func syncHandler(uid int64, bot *tgbotapi.BotAPI) {
	if cd, ok := users[uid]; ok {
		syncCal(cd)
		bot.Send(tgbotapi.NewMessage(uid, "Календарь синхронизирован."))
	} else {
		bot.Send(tgbotapi.NewMessage(uid, "Профиль не обнаружен."))
	}
}

func setTimeHandler(uid int64, bot *tgbotapi.BotAPI, args string) {
	if cd, ok := users[uid]; ok {
		if args == "" {
			bot.Send(tgbotapi.NewMessage(uid, "Текущая настройка: "+strconv.Itoa(cd.userTime)+" мин."))
			bot.Send(tgbotapi.NewMessage(uid, "Если хотите изменить время, то укажите его через: /settime <минуты>"))
			return
		}
		cd.userTime, _ = strconv.Atoi(args)
		UpdateElement(cd)
		bot.Send(tgbotapi.NewMessage(uid, "Время изменено."))
	} else {
		bot.Send(tgbotapi.NewMessage(uid, "Профиль не обнаружен."))
	}
}

func deleteHandler(uid int64, bot *tgbotapi.BotAPI) {
	if cd, ok := users[uid]; ok {
		DeleteElement(cd)
		delete(users, uid)
		bot.Send(tgbotapi.NewMessage(uid, "Профиль удален."))
	} else {
		bot.Send(tgbotapi.NewMessage(uid, "Профиль не обнаружен."))
	}
}
