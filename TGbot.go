package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"strings"
)

//var tickersMap = make(map[int][]string) //мапа внутри которой я запоминаю слайсы с тикерами
//тесты(как писать и как сделать свой код теситируемыми

var tickersMap = make(map[int]map[string]interface{}) //избавиться от глоб(структура со своим методом addticker)

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
			_, ok := tickersMap[UsID]
			if !ok {
				tickersMap[UsID] = make(map[string]interface{})
			}
			words := messageEnter(text)
			reply = addTicker(text, reply, UsID, words)

		case "mytickers":				//тикеры пользователя
			//reply = myTicker(reply)

		case "delete":
			//deleteTicker(text, reply)

		case "help":
			reply = ``

		default:
			reply = "Unknown command"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}

func start(reply string) string{
	reply = "Hello. I am your personal investment assistant. To find out what I can do, write "
	return reply
}

/*func myTicker(reply string, UsID int) string{//создать массив и сделать функцию(отдельную) сортировку тикеров в мапе
	var sls []string
	for _, v := range tickersMap[UsID]{
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

func addTicker(text string, reply string, UsID int, words []string) string{//сделать принятие сообщений отдельно от добавления
	for _, s := range words{
			for _, v := range tickersSlice {

				if v != s { //если такого тикера не существует
					reply = "Unknown command"

				}else if v == s { //если такой тикер найден
					tickersMap[UsID][s] = nil
					reply = "Ticker saved"
					break
				}
			}
	}
		/*}else {
			reply = "Ticker name not found"
		}
	}*/
	fmt.Println(tickersMap)
	return reply
}

func messageEnter(text string) []string{ //разбиваю ввод пользователя и удаляю
	words := strings.Fields(text)		//"/addticker", оставляя только названия тикеров
	words[0] = ""
	return words
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

//написать тест, разбить на функции
