package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"strings"
)

//var tickersMap = make(map[int][]string) // мапа внутри которой я запоминаю слайсы с тикерами

var tickersMap = make(map[int]map[string]interface{})


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

			UsID := update.Message.From.ID
			tickersMap[UsID] = make(map[string]interface{})
			reply = addTicker(text, reply, UsID)

		case "mytickers":				//тикеры пользователя
			//reply = myTicker(reply)

		case "delete":
			//deleteTicker(text, reply)

		case "help":
			reply = ""

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
/*
func myTicker(reply string, UsID int) string{
	var sls []string
	for _, v := range tickersMap{
		k := v
		sls = append(sls, k)
		return reply

	}
	return reply
}*/

/*
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
}*/

func addTicker(text string, reply string, UsID int) string{
	//var sls []string
	//var anMap = make(map[string]interface{})
	words := strings.Fields(text)
	for _, v := range words {
		if v != "/addticker" {
			myWord := v
			for _, v := range tickersSlice {
				if v != myWord { //если такого тикера не существует
					reply = "Unknown command"

				}else if v == myWord { //если такой тикер найден
					/*sls = append(sls, myWord)
					s := cap(sls)
					//dnd := tickersMap[UsID]
					for i, g := range sls {
						if i != s{
						anMap[g] = nil
						}
					}
					tickersMap[UsID] = anMap
					*/

				tickersMap[UsID][myWord] = nil
				//dnd := tickersMap[UsID][myWord]


				reply = "Ticker saved"
				break
				}
			}
		}else {
			reply = "Ticker name not found"
		}
	}
	fmt.Println(tickersMap)
	return reply
}

/*func deleteTicker(text string, reply string){
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
}*/

