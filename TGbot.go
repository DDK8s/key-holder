package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/xlab/closer"

	//"github.com/xlab/closer"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
)

var tickersSlice = []string{
	"ONE",
	"TWO",
	"THREE",
	"FOUR",
	"FIVE",
	"SIX",
	"SEVEN",
	"EIGHT",
	"NINE",
	"TEN",
}

func main(){
	var a tickersInt = &tickersStr{}
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

	a.diskReading(tickersMap)
	go a.autoSaving(tickersMap)

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
			reply = a.start(reply)

		case "addticker":
			a.mapValueChecker(tickersMap, UsID)
			words := a.fetchTickerName(text)
			reply = a.addTickers(reply, UsID, words, tickersMap)


		case "mytickers":
			tickers := a.sorting(tickersMap, UsID)
			reply = a.userTickers(reply, tickers)

		case "delete":
			words := a.fetchTickerName(text)
			reply = a.deleteTicker(reply, UsID, words, tickersMap)

		case "help":
			reply = a.help(reply)

		case "botoff":
			if UsID == 744515526 {
				a.dataSaving(tickersMap)
				reply = "Will see soon, Creator"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				bot.Send(msg)
				closer.Close()
			} else {
				reply = "You're not my Creator"
			}

		default:
			reply = "Unknown command. Write \"/help\" to see the list of commands."
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)

	}
	//обработка сигналов завершения работы
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer stop()

	select {
	case <-time.After(10 * time.Second):
		fmt.Println("missed signal")
	case <-ctx.Done():
		a.writeInJson
		stop()
		fmt.Println("signal received")
	}

}


//_____________________________Структуры_____________________________________________________________________________


type tickersInt interface {
	help(string) string
	start(string) string
	fetchTickerName(string) []string
	userTickers(string, []string) string
	autoSaving(map[int]map[string]interface{})
	writeInJson(map[int]map[string]interface{})
	diskReading(map[int]map[string]interface{})
	dataSaving(map[int]map[string]interface{})
	mapValueChecker(map[int]map[string]interface{}, int)
	sorting(map[int]map[string]interface{}, int) []string
	addTickers(string, int, []string, map[int]map[string]interface{}) string
	deleteTicker(string, int, []string, map[int]map[string]interface{}) string
}

type tickersStr struct {
	tickersMap map[int]map[string]interface{}
}


//_____________________________Методы________________________________________________________________________________


func (a *tickersStr) start(reply string) string{
	reply = "Hello. I am your personal investment assistant. To find out what I can do, write \"/help\"."
	return reply
}

func (a *tickersStr) help(reply string) string{
	reply = `◽Use the command "/addticker [ticker name]" to add a new ticker to your list of tickers.
◽Use the command "/delete [ticker name]" to remove the ticker from your list of tickers.
◽Use the command "/mytickers" to see a list of your tickers.`
	return reply
}

func (a *tickersStr) addTickers(reply string, UsID int, words []string, tickersMap map[int]map[string]interface{}) string{
	for _, s := range words{
		for _, v := range tickersSlice {
			if v != s { //если такого тикера не существует
				reply = "Unknown ticker."

			}else if v == s { //если такой тикер найден
				tickersMap[UsID][s] = nil
				reply = "Ticker saved."
				break
			}
		}
	}
	fmt.Println(tickersMap)
	return reply
}

func (a *tickersStr) deleteTicker(reply string, UsID int, words []string, tickersMap map[int]map[string]interface{}) string {
	for _, s := range words {
		for range tickersMap[UsID]{
			_, ok := tickersMap[UsID][s]
			if ok {
				delete(tickersMap[UsID], s)
			}
		}
	}
	fmt.Println(tickersMap[UsID])
	reply = "Ticker was deleted."
	return reply
}

func (a *tickersStr) sorting(tickersMap map[int]map[string]interface{}, UsID int) []string{
	tickers := make([]string, 0, len(tickersMap[UsID]))
	for v := range tickersMap[UsID] {
		tickers = append(tickers, v)
	}
	sort.Strings(tickers)
	return tickers
}

func (a *tickersStr) fetchTickerName(text string) []string{
	words := strings.Fields(text)
	words = words[1:]
	return words
}

func (a *tickersStr) writeInJson(tickersMap map[int]map[string]interface{}){

	file, _ := json.MarshalIndent(tickersMap, "", " ")
	_ = ioutil.WriteFile("test.json", file, 0644)

}

func (a *tickersStr) diskReading(tickersMap map[int]map[string]interface{}){
	file, _ := ioutil.ReadFile("test.json")
	json.Unmarshal(file, &tickersMap)
}

func (a *tickersStr) autoSaving(tickersMap map[int]map[string]interface{}) {
	if true {
		time.Sleep(5 * time.Minute)
		a.writeInJson(tickersMap)
	}

}

func (a *tickersStr) mapValueChecker(tickersMap map[int]map[string]interface{}, UsID int){
	_, ok := tickersMap[UsID]
	if !ok {
		tickersMap[UsID] = make(map[string]interface{})
	}
}

func (a *tickersStr) userTickers(reply string, tickers []string) string {
	if tickers == nil {
		reply = "Empty ticker list"
	}
	for _, v := range tickers{
		reply = reply + v + " "
	}
	return reply
}

func (a *tickersStr) dataSaving(tickersMap map[int]map[string]interface{}) {
	go a.writeInJson(tickersMap)
	var wg sync.WaitGroup
	wg.Wait()
}

//___________________________________________________________________________________________________________________

