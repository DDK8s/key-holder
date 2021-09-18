package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"go/types"
	"log"
	"strings"
)

var myTickers []string//массив тикеров пользователя, который пользователь может редактировать

var tickersMap = make(map[int]types.Slice) // мапа внутри которой я запоминаю слайсы с тикерами

var tickersSlice = []string{
	"ONE",
	"TWO",
	"THREE",
	"FOUR",
}

func main(){

	bot, err := tgbotapi.NewBotAPI("1935733666:AAGj-bDMkUR6DZIqwiNjhDJCbomieEkVZYo")
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

	// Обновления канала
	for update := range updates {
		text := update.Message.Text	//текст сообщения
		var reply string	//ответ на сообщение
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)


		switch update.Message.Command() {
			case "start":
				reply = start(reply)

			case "addticker":

				reply = addTicker(text, reply)
				tickersMap[update.Message.From.ID] = myTickers// я не реализовал это внутри функции,
				myTickers = nil								 //т.к. не могу передать update.Message.From.ID в неё
				
				fmt.Println(tickersMap)//просто для удобства проверки содержимого
				
				
				/*идея заключается в том, чтобы запомнить тикеры myTickers пользователя, передать в мапу по уникальному
				update.Message.From.ID, сбросить слайс myTickers и по кругу с каждым.
				Проблема ещё в том, что на 1 ключ мапы приходится 1 элемент и думал над чем-то вроде:
				for i, v := range myTickers {
				tickersMap[update.Message.From.ID] = tickersMap[update.Message.From.ID] + v
				}, но получается каша
				 */



			case "mytickers":				//тикеры пользователя
				reply = myTicker(reply)

			case "delete":
				deleteTicker(text, reply)

			default:
				reply = "Unknown command"
			}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}

func start(reply string) string{
	reply = "Hello, i'm a telegram bot"
	return reply
}

func myTicker(reply string) string{
	var myTickersSlice []string
	if myTickers != nil {
		for _, l := range myTickers{
			myTickersSlice = append(myTickersSlice, l)
			s := len(myTickersSlice)
			for i, _ := range myTickersSlice{
				if i != s {
					reply = strings.Join(myTickersSlice, " ")
				}
			}
		}
	}else{
		reply = "Empty ticker list"
	}
	return reply
}

func addTicker(text string, reply string) string{

	words := strings.Fields(text)
	for _, v := range words {
		if v != "/addticker" {
			myWord := v

			for _, v := range tickersSlice {
				if v != myWord { //если такого тикера не существует
					reply = "Unknown command"

				}else if v == myWord { //если такой тикер найден
					//for range
					myTickers = append(myTickers, myWord)

					reply = "Ticker saved"
					break
				}
			}
		}
	}
	fmt.Println(myTickers)
	return reply
}

func deleteTicker(text string, reply string){
	words := strings.Fields(text)
	for _, t := range words {
		if t != "/addticker" {
			myWord := t

			for i, v := range tickersSlice {
				if v != myWord { //если такого тикера не существует
					reply = "Ticker not found"

				} else if v == myWord { //если такой тикер найден
					myTickers = append(myTickers[:i], myTickers[i+1])
					reply = "Ticker deleted"
					break
				}
			}
		}
	}
	fmt.Println(myTickers)
}
