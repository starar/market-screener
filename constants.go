package main

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	BotToken   = "1268262051:AAEaznI5Uc45PGxmBkSXBEwHcIyro0UQ5CI"
	WebhookURL = "https://screener-project.herokuapp.com"
	//WebhookURL = "https://30535d304af1.ngrok.io"
)

const (
	tooLong               = 50
	cycle   time.Duration = 30 * time.Minute
)

const (
	waitingStart = iota
	waitingAnything
	waitingListName
	waitingCompanyName
	waitingAcception
	waitingMode
	waitingLowerBound
	waitingUpperBound
	waitingListPicking
	waitingManage
)

const (
	listActive = iota
	listInactive
	listDeleted
	listPrepare
	modePrice
	modeCapital
	modePercent
)

var modeButtons = [][]tgbotapi.KeyboardButton{NewButton(priceButton), NewButton(capitalButton)}
var manageButtons = [][]tgbotapi.KeyboardButton{NewButton(deleteButtom), NewButton(cancelButtom)}

var (
	startCommand  = "/start"
	newCommand    = "/new"
	listCommand   = "/list"
	stopCommand   = "/stop"
	finishCommand = "/finish"
	manageCommand = "/manage"
)

var (
	notFoundButtom   = `Нужной компании нет в списке`
	priceButton      = `По цене за акцию`
	capitalButton    = `По рыночной капитализации`
	activateButton   = `Активировать список`
	deactivateButton = `Остановить отслеживание`
	deleteButtom     = `Удалить список`
	cancelButtom     = `Отмена`
)

var (
	helpMessage = `Это бот-скринер акций компаний Московской биржи. 

Вы можете создавать портфели акций, выставляя интересные границы для стоимости. Если цена войдёт в эти границы, то бот автоматически пришлёт уведомление. Так, вы можете почти не тратя время на анализ рынка, делать верные решения в инвестировании.

С помощью /manage портфели можно просматривать, останавливать их отслеживание ботом, удалять.
Давайте создадим портфель командой /new.`
	newListMessage    = `Задайте название портфеля`
	badListMessage    = `Название занято. Введите другое.`
	newCompanyMessage = `Введите название компании.
Например: сбер, yandex, Нефтегаз.`
	chooseModeMessage   = `Выберите тип фильтра:`
	notFoundMessage     = `Ничего не найдено. Уточните название компании.`
	clarifyMessage      = `Уточните и введите запрос снова`
	alreadyExistMessage = `Компания уже добавлена. Введите другую.`
	selectMessage       = `Выберите компанию из списка:`
	companyInfoMessage  = `Тикер: %s
Имя комапнии: %s
Цена за акцию: %.2f₽
Рыночная капитализация: %s`
	tryAgainMessage    = `Попробуйте ещё раз`
	lowerBoundMessage  = `Отлично! Теперь введите стоимость интересной нижней границы`
	badNumbaerMessage  = `Введите неотрицательное число`
	compareMessage     = `Верхняя граница должна быть больше нижней`
	upperBoundMessage  = `Теперь введите стоимость верхней границы`
	finalTickerMessage = `Для тикера %s %s от %.2f до %.2f
Введите название следующей компании, чтобы добавить её в портфель. 
Или /finish, чтобы завершить создание портфеля.`
	capitalMessage     = `выбран объём капитализации`
	priceMessage       = `выбрана цена акции`
	finishMessage      = `Портфель успешно создан! Вызовите /manage, чтобы просмотреть свои портфели и управлять ими.`
	lineListMessage    = `#%d %s — %s`
	isActiveMessage    = `активен`
	isInactiveMessage  = `остановлен`
	pickListMessage    = `Выберите портфель, чтобы просмотреть его:`
	tickerPriceMessage = `%s — цена от %.2f до %.2f
`
	tickerCapitalMessage = `%s — объём капитализации от %.2f до %.2f
`
	empltyListMessage       = `Список пуст`
	successMessage          = `Хорошо!`
	misunderstandingMessage = `Кажется, я вас не понял. Попробуйте /start.`
	stopMessage             = `Бот остановлен. Чтобы запустить бота, нажмите /start.`
	matchMessage            = `Сработало!
`
	infoMessage = `В портфеле %s в границы входят компании: 
`
	tickerMessage = `%s (%.2f₽)  `
)
