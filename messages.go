package main

import (
	"database/sql"
	"fmt"
	"market-screener/moex"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Result struct {
	Chatid int64
	Text   string
}

type ListItem struct {
	userid int
	name   string
}

func NewButton(s string) []tgbotapi.KeyboardButton {
	return []tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton(s)}
}

func Send(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(
		chatID,
		text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	bot.Send(msg)
}

func SendKeyboard(bot *tgbotapi.BotAPI, chatID int64, text string, keys [][]tgbotapi.KeyboardButton) {
	keyboard := tgbotapi.ReplyKeyboardMarkup{
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
		Keyboard:        keys,
	}

	msg := tgbotapi.NewMessage(
		chatID,
		text)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func AutoUpdate(db *sql.DB, ch chan<- Result) {
	ticker := time.Tick(cycle)

	moex.UpdateSecurities(db)

	for range ticker {
		moex.UpdateSecurities(db)
		CheckMatch(db, ch, 0)
	}
}

func CheckMatch(db *sql.DB, ch chan<- Result, barrier int) {
	prices := moex.GetTickersAll(db)
	tickerMap := make(map[string]moex.Security)
	for _, sec := range prices {
		tickerMap[sec.Ticker] = sec
	}

	items := moex.GetItemsAll(db)
	listids, userids, names := moex.GetActiveLists(db, listActive)
	listMap := make(map[int]ListItem)
	usersMap := make(map[int]map[int][]string)
	for i, listid := range listids {
		listMap[listid] = ListItem{userids[i], names[i]}
	}
	fmt.Println(listMap)

	for _, item := range items {
		var listIt ListItem
		var flag bool
		if listIt, flag = listMap[item.Listid]; flag == false {
			continue
		}
		flag = false
		if item.Mode == modePrice {
			if item.Lower <= tickerMap[item.Ticker].Price && tickerMap[item.Ticker].Price <= item.Upper {
				flag = true
			}
		} else {
			if item.Lower <= tickerMap[item.Ticker].Capital && tickerMap[item.Ticker].Capital <= item.Upper {
				flag = true
			}
		}
		if flag {
			if len(usersMap[listIt.userid]) == 0 {
				usersMap[listIt.userid] = make(map[int][]string)
			}
			usersMap[listIt.userid][item.Listid] = append(usersMap[listIt.userid][item.Listid], item.Ticker)
		}
	}

	fmt.Println(usersMap)

	for userid, lists := range usersMap {
		if len(listMap) == 0 {
			continue
		}
		var s strings.Builder
		s.WriteString(matchMessage)

		for listid, tickers := range lists {
			if len(tickers) == 0 {
				continue
			}
			s.WriteString(fmt.Sprintf(infoMessage, listMap[listid].name))
			for _, ticker := range tickers {
				s.WriteString(fmt.Sprintf(tickerMessage, tickerMap[ticker].ShortName, tickerMap[ticker].Price))
			}
			s.WriteString("\n")
		}
		if barrier != 0 && userid != barrier {
			continue
		}
		ch <- Result{int64(userid), s.String()}
	}
}
