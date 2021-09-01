package main

import (
	"context"
	//"encoding/json"
//	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"log"
	//"os"

	"time"
)

const (
	tgtoken   = "t.nbVvVcHGxEFof8uIWo7uywaAubEo2aULyhzPLXjuesSe9_Y8TzQUar37G_Hk1fsaY1xQMnrfFf0DGzk4Nj3R5g"
	timeout = time.Second * 3
	url     = "https://api-invest.tinkoff.ru/openapi/sandbox/sandbox/register"
)

func main() {
	ReturnMessage()
}

func ReturnMessage() {
	bot, err := tgbotapi.NewBotAPI(tgtoken)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}
	// В канал updates будут приходить все новые сообщения.
	for update := range updates {
		// Создав структуру - можно её отправить обратно боту
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}

func ticker() { //cм материал 3
	client := sdk.NewRestClient(tgtoken)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//получение инструмента по тикеру
	instruments, err := client.InstrumentByTicker(ctx, "TCS")
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v\n", instruments)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}


/*
1.	Создать бота+
2.	Библиотеки+
2.1.	Библеотека для ТГбота+
2.2.	Библиотека для streaming OpenAPI+
3.	Связаться со Stream OpenAPI+
4.	Получить котировки
4.1.	Вытащить данные из OpenAPI
4.2.	Получение изменений в реальном времени
4.3.	Запомнить изменения
5.	Создать алгоритм вычисления точки ema за n дней
6.	Вычислить точку ema
6.1.	За последние n дней
6.2.	За указанный период
7.	Создать команды у бота
8.	Привязать команды к нашему go

____Материалы___________________________________________________

1. Тинькофф Streaming OpenAPI:
https://tinkoffcreditsystems.github.io/invest-openapi/marketdata/

2. Код для инфестиций:
https://github.com/TinkoffCreditSystems/invest-openapi-go-sdk/blob/master/examples/main.go


3. sdk tinkoff api
//стр. 169 https://github.com/TinkoffCreditSystems/invest-openapi-go-sdk/blob/master/examples/main.go#L169
 */
