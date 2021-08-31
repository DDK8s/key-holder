package main

import (
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"time"
)

func telegramBot() {
	//подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI("1935733666:AAGj-bDMkUR6DZIqwiNjhDJCbomieEkVZYo")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Aythorized on account %s", bot.Self.UserName)

	//Устанавливаем время обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//получаем обновления от бота
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
	}

}



/*
1.	Создать бота+
2.	Библиотеки+
2.1.	Библеотека для ТГбота+
2.2.	Библиотека для streaming OpenAPI
3.	Связаться со Stream OpenAPI
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

Тинькофф Streaming OpenAPI:
https://tinkoffcreditsystems.github.io/invest-openapi/marketdata/
