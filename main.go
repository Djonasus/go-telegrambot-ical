package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/emersion/go-ical"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/slices"
)

var (
	users []CalData
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
			//continue
			//if update.ChatMember.NewChatMember.HasLeft() {
			//if update.MyChatMember.NewChatMember.HasLeft() {
			if update.MyChatMember.NewChatMember.Status == "left" || update.MyChatMember.NewChatMember.Status == "kicked" {
				cd, i := FindUserById(update.MyChatMember.Chat.ID, users)
				if cd.userID != 0 {
					e := os.Remove(cd.userCalendar)
					if e != nil {
						log.Fatal(e)
					}
					fmt.Println("User " + strconv.FormatInt(update.MyChatMember.Chat.ID, 10) + " gone(")
					DeleteElement(&cd)
					users[i] = CalData{}
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
			cd, i := FindUserById(update.FromChat().ID, users)

			if cd.userID == 0 {
				bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "У вас не создан профиль. Введите /start!"))
				continue
			}

			switch cd.userState {
			case "create", "update":
				err := DownloadFile("calendars/"+strconv.FormatInt(update.FromChat().ID, 10)+".ical", update.Message.Text)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.FromChat().ID, err.Error()))
					continue
				}
				users[i].userState = "listen"
				users[i].userCalendar = "calendars/" + strconv.FormatInt(update.FromChat().ID, 10) + ".ical"
				users[i].userURL = update.Message.Text

				err = UpdateElement(&users[i])
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.FromChat().ID, err.Error()))
				}

				bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Скачивание календаря завершено!"))
				//bot.Send(tgbotapi.NewMessage(update.FromChat().ID, err.Error()))
			}
		}
	}
}

// GOROUTINES

func syncCals() {
	for {
		if len(users) != 0 {
			for _, cd := range users {
				go syncCal(&cd)
			}
			fmt.Println("All cals synced")
		}
		time.Sleep(10 * time.Minute)
	}
}

func syncCal(cd *CalData) {
	if cd.userID == 0 || cd.userURL == "" || cd.userState == "create" || cd.userState == "update" {
		//fmt.Println("EMPTY!")
		return
	}
	err := DownloadFile(cd.userCalendar, cd.userURL)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func callMe(bot *tgbotapi.BotAPI) {
	loc, _ := time.LoadLocation("")
	//call logic
	for {
		if len(users) != 0 {
			for i, cd := range users {
				if cd.userCalendar == "" {
					continue
				}
				eve, err := getEventsNames(cd.userCalendar)
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, ev := range eve {
					nme, _ := ev.Props.Text(ical.PropSummary)
					startTm, _ := ev.DateTimeStart(loc)
					dif := startTm.Sub(time.Now()).Minutes()

					if dif <= float64(cd.userTime) && dif > 0 && !slices.Contains(cd.userShowedEvents, nme) {
						bot.Send(tgbotapi.NewMessage(cd.userID, "Событие "+nme+" скоро начнется! Осталось "+fmt.Sprintf("%.0f", dif)+" минут!"))
						users[i].userShowedEvents = append(users[i].userShowedEvents, nme)
						UpdateElement(&users[i])
					}
				}
			}
		}
		time.Sleep(time.Minute)
	}
}

/*func main() {

	ok := CalData{userID: 160, userCalendar: "aboba", userURL: "beep", userState: "start", userShowedEvents: []string{"ok", "neok"}}

	//NewElement(ok)

	//LoadData()
	//for _, v := range users {
	//fmt.Println(v.userID, v.userCalendar, v.userURL)
	//}

	ok.userURL = "Not Aboba!"
	UpdateElement(&ok)

	LoadData()
	for _, v := range users {
		fmt.Println(v.userID, v.userCalendar, v.userURL, v.userShowedEvents)
	}
}
*/
