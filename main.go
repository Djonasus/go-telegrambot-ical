package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/slices"
)

var (
	users map[int64]*CalData
)

func main() {

	LoadData()

	if len(os.Args) == 1 || os.Args[1] == "" {
		panic("Ошибка! Укажите токен бота ./bot <ВАШ_ТОКЕН>")
	}

	bot, err := tgbotapi.NewBotAPI(os.Args[1])
	//bot, err := tgbotapi.NewBotAPI("6676298340:AAHkSzB-EpE_Povq86_N38EKu4lUAcRj7pM")
	if err != nil {
		panic(err)
	}
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	go callMe(bot)
	go syncCals()

	for update := range updates {
		if update.Message == nil {
			fmt.Println("Empty Update") // ignore any non-Message updates
			if update.MyChatMember.NewChatMember.Status == "left" || update.MyChatMember.NewChatMember.Status == "kicked" {
				if cd, ok := users[update.MyChatMember.Chat.ID]; ok {
					fmt.Println("User " + strconv.FormatInt(update.MyChatMember.Chat.ID, 10) + " gone(")
					DeleteElement(cd)
					delete(users, update.MyChatMember.Chat.ID)
				}
			}
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "saymyname":
				bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Твое имя:"))
				bot.Send(tgbotapi.NewMessage(update.FromChat().ID, update.FromChat().UserName))
			case "start":
				startHandler(update.FromChat().ID, bot)
			case "debug":
				debugHandler(update.FromChat().ID, bot)
			case "events":
				eventsHandler(update.FromChat().ID, bot)
			case "sync":
				syncHandler(update.FromChat().ID, bot)
			case "settime":
				setTimeHandler(update.FromChat().ID, bot, update.Message.CommandArguments())
			case "delete":
				deleteHandler(update.FromChat().ID, bot)
			default:
				bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Неизвестная команда"))
			}
		} else {
			//FOR URL OR SOMETHING ELSE

			if cd, ok := users[update.FromChat().ID]; ok {
				switch cd.userState {
				case "create", "update":
					cd.userURL = update.Message.Text
					events, err := GetEvents(*cd)
					if err != nil {
						bot.Send(tgbotapi.NewMessage(update.FromChat().ID, err.Error()))
						continue
					}
					cd.userEvents = events
					cd.userState = "listen"

					err = UpdateElement(cd)
					if err != nil {
						bot.Send(tgbotapi.NewMessage(update.FromChat().ID, err.Error()))
					}

					bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Скачивание календаря завершено!"))
				}
			} else {
				bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "У вас не создан профиль. Введите /start!"))
				continue
			}
		}
	}
}

// GOROUTINES

func syncCals() {
	for {
		if len(users) != 0 {
			for _, v := range users {
				go syncCal(v)
			}
			fmt.Println("All cals synced")
		}
		time.Sleep(10 * time.Minute)
	}
}

func syncCal(cd *CalData) {
	loc, _ := time.LoadLocation("Local")
	if cd.userID == 0 || cd.userURL == "" || cd.userState == "create" || cd.userState == "update" {
		return
	}

	newEvents, err := GetEvents(*cd)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, v := range newEvents {
		if !slices.Contains(cd.userEvents, v) {
			cd.userEvents = append(cd.userEvents, v)
		}
	}
	for i := len(cd.userEvents) - 1; i >= 0; i-- {
		dif := cd.userEvents[i].Date.Sub(time.Now().In(loc)).Minutes()
		if dif <= 0 {
			cd.userEvents = removeEvent(cd.userEvents, i)
		}
	}
	UpdateElement(cd)
}

func callMe(bot *tgbotapi.BotAPI) {
	loc, _ := time.LoadLocation("Local")
	//call logic
	for {
		if len(users) != 0 {
			for _, cd := range users {
				fmt.Println("Call user: " + strconv.FormatInt(cd.userID, 10))

				if cd.userEvents == nil {
					continue
				}
				for ei, _ := range cd.userEvents {

					dif := cd.userEvents[ei].Date.Sub(time.Now().In(loc)).Minutes()
					fmt.Println(time.Now(), cd.userEvents[ei])
					fmt.Println(dif)
					if cd.userEvents[ei].Showed == false && dif <= float64(cd.userTime) {
						fmt.Println("Send Call: " + cd.userEvents[ei].Name)
						bot.Send(tgbotapi.NewMessage(cd.userID, "Событие "+cd.userEvents[ei].Name+" скоро начнется! Осталось "+fmt.Sprintf("%.0f", dif)+" мин."))
						cd.userEvents[ei].Showed = true
						UpdateElement(cd)
					}
				}
			}
		}
		time.Sleep(time.Minute)
	}
}
