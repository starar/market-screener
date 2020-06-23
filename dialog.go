package main

import (
	"database/sql"
	"fmt"
	"market-screener/moex"
	"math"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func processUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *sql.DB) {
	if update.Message == nil {
		return
	}

	state, listidSt, tickerSt, err := moex.GetState(db, update.Message.From.ID)
	if err != nil {
		fmt.Println(err)
		state = waitingStart
	}
	fmt.Println("state ", state)

	message := update.Message.Text
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	if state == waitingStart {
		if message == startCommand {
			Send(bot, chatID, helpMessage)
			moex.UpdateListsAll(db, userID, listActive)
			moex.SaveState(db, userID, listidSt, tickerSt, waitingAnything)
			return
		}
	}

	if state == waitingAnything {
		switch message {
		case startCommand:
			Send(bot, chatID, helpMessage)
			return

		case newCommand:
			Send(bot, chatID, newListMessage)
			moex.SaveState(db, userID, listidSt, tickerSt, waitingListName)
			return

		case manageCommand:
			ids, names, states := moex.GetListsAll(db, userID, listDeleted)
			var buttons [][]tgbotapi.KeyboardButton
			for i, name := range names {
				var line string
				if states[i] == listActive {
					line += fmt.Sprintf(lineListMessage, ids[i], name, isActiveMessage)
				} else {
					line += fmt.Sprintf(lineListMessage, ids[i], name, isInactiveMessage)
				}
				buttons = append(buttons, NewButton(line))
			}
			if len(buttons) == 0 {
				Send(bot, chatID, empltyListMessage)
				return
			}
			SendKeyboard(bot, chatID, pickListMessage, buttons)
			moex.SaveState(db, userID, listidSt, tickerSt, waitingListPicking)
			return

		case stopCommand:
			moex.UpdateListsAll(db, userID, listInactive)
			Send(bot, chatID, stopMessage)
			moex.SaveState(db, userID, listidSt, tickerSt, waitingStart)
			return
		}
	}

	if state == waitingListName {
		if moex.IsListNameFree(db, userID, message) == false {
			Send(bot, chatID, badListMessage)
			return
		}

		moex.InsertList(db, userID, message, listPrepare)
		if listidSt, err = moex.GetListByName(db, userID, message); err != nil {
			fmt.Println("getlistid:", err)
		}

		Send(bot, chatID, newCompanyMessage)
		moex.SaveState(db, userID, listidSt, tickerSt, waitingCompanyName)
		return
	}

	if state == waitingCompanyName {
		if message == finishCommand {
			Send(bot, chatID, finishMessage)
			moex.UpdateList(db, userID, listidSt, listActive)
			CheckMatch(db, notifications, userID)
			moex.SaveState(db, userID, listidSt, tickerSt, waitingAnything)
			return
		}

		response := moex.FindCompany(db, message)
		if len(response) == 0 {
			Send(bot, chatID, notFoundMessage)
			return
		}
		if len(response) > tooLong {
			response = response[:tooLong]
		}

		var buttons [][]tgbotapi.KeyboardButton
		buttons = append(buttons, NewButton(notFoundButtom))
		for _, sec := range response {
			buttons = append(buttons, NewButton(sec.Ticker+"\t"+sec.Name))
		}
		SendKeyboard(bot, chatID, selectMessage, buttons)

		moex.SaveState(db, userID, listidSt, tickerSt, waitingAcception)
		return
	}

	if state == waitingAcception {
		name := strings.Split(message, " ")
		tickerSt = name[0]

		if message == notFoundButtom {
			Send(bot, chatID, clarifyMessage)
			moex.SaveState(db, userID, listidSt, tickerSt, waitingCompanyName)
			return
		}
		if moex.IsTickerExist(db, tickerSt) == false {
			Send(bot, chatID, notFoundMessage)
			moex.SaveState(db, userID, listidSt, tickerSt, waitingCompanyName)
			return
		}
		if moex.IsItemTickerFree(db, listidSt, tickerSt) == false {
			Send(bot, chatID, alreadyExistMessage)
			moex.SaveState(db, userID, listidSt, tickerSt, waitingCompanyName)
			return
		}

		price, mc := moex.GetTicker(db, tickerSt)
		Send(bot, chatID, fmt.Sprintf(companyInfoMessage, tickerSt,
			strings.Join(name[1:], " "), price, prettyFormat(mc)))
		SendKeyboard(bot, chatID, chooseModeMessage, modeButtons)
		moex.SaveState(db, userID, listidSt, tickerSt, waitingMode)
		return
	}

	if state == waitingMode {
		switch message {
		case priceButton:
			moex.InsertItem(db, listidSt, tickerSt, modePrice)
		case capitalButton:
			moex.InsertItem(db, listidSt, tickerSt, modeCapital)
		default:
			Send(bot, chatID, tryAgainMessage)
			SendKeyboard(bot, chatID, chooseModeMessage, modeButtons)
			return
		}

		Send(bot, chatID, lowerBoundMessage)
		moex.SaveState(db, userID, listidSt, tickerSt, waitingLowerBound)
		return
	}

	if state == waitingLowerBound {
		lowerBound, err := strconv.ParseFloat(message, 64)
		if err != nil || lowerBound < 0 {
			Send(bot, chatID, badNumbaerMessage)
			return
		}

		moex.UpdateItem(db, listidSt, tickerSt, lowerBound, 0)
		Send(bot, chatID, upperBoundMessage)
		moex.SaveState(db, userID, listidSt, tickerSt, waitingUpperBound)
		return
	}

	if state == waitingUpperBound {
		upperBound, err := strconv.ParseFloat(message, 64)
		mode, lowerBound, _, _ := moex.GetItem(db, listidSt, tickerSt)
		if err != nil || lowerBound < 0 {
			Send(bot, chatID, badNumbaerMessage)
			return
		}
		if lowerBound >= upperBound {
			Send(bot, chatID, compareMessage)
			return
		}
		moex.UpdateItem(db, listidSt, tickerSt, lowerBound, upperBound)

		if mode == modeCapital {
			Send(bot, chatID, fmt.Sprintf(finalTickerMessage, tickerSt,
				capitalMessage, lowerBound, upperBound))
		} else {
			Send(bot, chatID, fmt.Sprintf(finalTickerMessage, tickerSt,
				priceMessage, lowerBound, upperBound))
		}

		moex.SaveState(db, userID, listidSt, tickerSt, waitingCompanyName)
		return
	}

	if state == waitingListPicking {
		listidSt, err = strconv.Atoi(strings.Split(message, " ")[0][1:])
		if err != nil {
			Send(bot, chatID, tryAgainMessage)
			moex.SaveState(db, userID, listidSt, tickerSt, waitingAnything)
			return
		}
		tickers, modes, lower, upper := moex.GetItemsByList(db, listidSt)
		if len(tickers) == 0 {
			Send(bot, chatID, empltyListMessage)
			moex.UpdateList(db, userID, listActive, listDeleted)
			moex.SaveState(db, userID, listidSt, tickerSt, waitingAnything)
			return
		}
		var s strings.Builder
		for i, ticker := range tickers {
			if modes[i] == modePrice {
				fmt.Fprintf(&s, tickerPriceMessage, ticker, lower[i], upper[i])
			} else {
				fmt.Fprintf(&s, tickerCapitalMessage, ticker, lower[i], upper[i])
			}
		}

		listState, _ := moex.GetListState(db, userID, listidSt)
		if listState == listDeleted {
			moex.SaveState(db, userID, listidSt, tickerSt, waitingAnything)
			return
		}

		var buttons [][]tgbotapi.KeyboardButton
		if listState == listActive {
			buttons = append(buttons, NewButton(deactivateButton))
		} else {
			buttons = append(buttons, NewButton(activateButton))
		}
		buttons = append(buttons, manageButtons...)

		SendKeyboard(bot, chatID, s.String(), buttons)
		moex.SaveState(db, userID, listidSt, tickerSt, waitingManage)
		return
	}

	if state == waitingManage {
		switch message {
		case activateButton:
			moex.UpdateList(db, userID, listidSt, listActive)
		case deactivateButton:
			moex.UpdateList(db, userID, listidSt, listInactive)
		case deleteButtom:
			moex.UpdateList(db, userID, listidSt, listDeleted)
		}
		Send(bot, chatID, successMessage)
		moex.SaveState(db, userID, listidSt, tickerSt, waitingAnything)
		return
	}

	Send(bot, chatID, misunderstandingMessage)
	return
}

func prettyFormat(f float64) string {
	if math.Pow10(5) <= f && f < math.Pow10(8) {
		return fmt.Sprintf("%.f тыс.", f/math.Pow10(3))
	}
	if math.Pow10(8) <= f && f < math.Pow10(11) {
		return fmt.Sprintf("%.f млн.", f/math.Pow10(6))
	}
	if f >= math.Pow10(11) {
		return fmt.Sprintf("%.f млрд.", f/math.Pow10(9))
	}
	return fmt.Sprintf("%.2f", f)
}
