package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"market-screener/moex"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var notifications chan Result = make(chan Result, 1)

func main() {
	var store *sql.DB = moex.CreateTable()
	defer store.Close()
	go AutoUpdate(store, notifications)

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		panic(err)
	}

	updates := bot.ListenForWebhook("/")

	port := os.Getenv("PORT")
	go http.ListenAndServe(":"+port, nil)
	fmt.Println("Start listen", port)

	for {
		select {
		case update := <-updates:
			processUpdate(bot, update, store)
		case notification := <-notifications:
			Send(bot, notification.Chatid, notification.Text)
		}
	}
}
