package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"strings"
)

//тесты(как писать и как сделать свой код теситируемыми

//var tickersMap = make(map[int]map[string]interface{}) //избавиться от глоб(структура со своим методом addticker)

var tickersSlice = []string{
	"ONE",
	"TWO",
	"THREE",
	"FOUR",
}

func main(){
	var tickersMap = make(map[int]map[string]interface{})
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

		text := update.Message.Text
		var reply string

		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		UsID := update.Message.From.ID

		switch update.Message.Command() {
		case "start":
			reply = start(reply)

		case "addticker":
			_, ok := tickersMap[UsID]
			if !ok {
				tickersMap[UsID] = make(map[string]interface{})
			}
			words := tickerNamePulling(text)
			reply = addTicker(reply, UsID, words, tickersMap)

		case "mytickers":				//тикеры пользователя
			reply = userTickers(reply, UsID, tickersMap)

		case "delete":
			words := tickerNamePulling(text)
			deleteTicker(reply, UsID, words, tickersMap)

		case "help":
			reply = `◽Use the command "/addticker [ticker name]" to add a new ticker to your list of tickers.
◽Use the command "/delete [ticker name]" to to remove the ticker from your list of tickers.
◽Use the command "/mytickers" to see a list of your tickers.`

		default:
			reply = "Unknown command"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}

func start(reply string) string{
	reply = "Hello. I am your personal investment assistant. To find out what I can do, write \"/help\"."
	return reply
}

func userTickers(reply string, UsID int, tickersMap map[int]map[string]interface{}) string {//создать массив и сделать функцию(отдельную) сортировку тикеров в мапе
	for i := range tickersMap[UsID] {
		reply = reply + i + " "
	}
	return reply
}

func tickerNamePulling(text string) []string{ //разбиваю ввод пользователя и удаляю
	words := strings.Fields(text)		   	//"/addticker", оставляя только названия тикеров
	words = words[1:]
	return words
}

func addTicker(reply string, UsID int, words []string, tickersMap map[int]map[string]interface{}) string{
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
	fmt.Println(tickersMap)
	return reply
}

func deleteTicker(reply string, UsID int, words []string, tickersMap map[int]map[string]interface{}) string {
	for _, s := range words {
		for range tickersMap[UsID]{
			_, ok := tickersMap[UsID][s]
			if ok {
				delete(tickersMap[UsID], s)

			}
		}
	}
	fmt.Println(tickersMap[UsID])
	reply = "Ticker was deleted"
	return reply
}

//написать тест, разбить на функции
