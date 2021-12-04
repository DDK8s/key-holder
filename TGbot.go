package main

import (

	"awesomeProject/5/Storage"
	"flag"
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"

	//"awesomeProject/5/Buisnes"
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
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

var token = flag.String("token", "", "t.-6-fOCXjCkzJGhM2QLtsOvysQdLGiuWUUq5T0VJCVJkYpEbg1N8PpPHSdttnhXeKPeKSG2j1byMb-i4rbfg6Kg")
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

func main() {

	var a tickersInt = &tickersStr{}
	var tickersVault = make(map[int]map[string]interface{})
	bot, err := tgbotapi.NewBotAPI("1935733666:AAGh55XNSeps0LYMpbAiAM-zKsZT-EQEFKc")
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

	Storage.ReadFromJson(tickersVault)
	go Storage.AutoSave(tickersVault)

	//поместил сюда L113-L126, ловлю сигнал остановки сразу же

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
			Storage.CheckingMapInitialization(tickersVault, UsID)
			words := a.fetchTickerName(text)
			reply = a.addTickers(reply, UsID, words, tickersVault)

		case "mytickers":
			tickers := a.sorting(tickersVault, UsID)

			reply = Storage.ReturnUserTickerList(reply,tickers)

		case "delete":
			words := a.fetchTickerName(text)
			reply = a.deleteTicker(reply, UsID, words, tickersVault)

		case "help":
			reply = a.help(reply)

		case "botoff":
			if UsID == 744515526 {
				Storage.DataSave(tickersVault)
				reply = "Will see soon, Creator"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				bot.Send(msg)
				os.Exit(1)
			} else {
				reply = "You're not my Creator"
			}

		default:
			reply = "Unknown command. Write \"/help\" to see the list of commands."
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)

	}
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	stop()
	select {
	case <-time.After(10 * time.Second):
		fmt.Println("missed signal")
		//документация контекста
	case <-ctx.Done():
		Storage.WriteInJson(tickersVault)
		fmt.Println("signal received")
	}
}

//_____________________________Структуры_____________________________________________________________________________

type tickersInt interface {
	help(string) string
	start(string) string
	fetchTickerName(string) []string
	ReturnUserTickerList(string, []string) string
	sorting(map[int]map[string]interface{}, int) []string
	addTickers(string, int, []string, map[int]map[string]interface{}) string
	deleteTicker(string, int, []string, map[int]map[string]interface{}) string
}

type tickersStr struct {
	tickersVault map[int]map[string]interface{}
}


//_____________________________Методы________________________________________________________________________________


func (a *tickersStr) start(reply string) string {
	reply = "Hello. I am your personal investment assistant. To find out what I can do, write \"/help\"."
	return reply
}

func (a *tickersStr) help(reply string) string {
	reply = `◽Use the command "/addticker [ticker name]" to add a new ticker to your list of tickers.
◽Use the command "/delete [ticker name]" to remove the ticker from your list of tickers.
◽Use the command "/mytickers" to see a list of your tickers.`
	return reply
}

func (a *tickersStr) addTickers(reply string, UsID int, words []string, tickersVault map[int]map[string]interface{}) string {
	for _, s := range words {
		for _, v := range tickersSlice {
			if v != s { //если такого тикера не существует
				reply = "Unknown ticker name " + s + "."

			} else if v == s { //если такой тикер найден
				tickersVault[UsID][s] = nil
				reply = "Ticker " + s + " saved."
				break
			}
		}
	}
	fmt.Println(tickersVault)
	return reply
}

func (a *tickersStr) deleteTicker(reply string, UsID int, words []string, tickersVault map[int]map[string]interface{}) string {
	for _, s := range words {
		_, ok := tickersVault[UsID][s]
		if ok {
			delete(tickersVault[UsID], s)
		}
	}
	fmt.Println(tickersVault[UsID])
	reply = "Ticker was deleted."
	return reply
}

//попробовать отсортировать в мапе
//мапы не имеют методов сортировки
func (a *tickersStr) sorting(tickersVault map[int]map[string]interface{}, UsID int) []string {
	tickers := make([]string, 0, len(tickersVault[UsID]))
	for v := range tickersVault[UsID] {
		tickers = append(tickers, v)
	}
	sort.Strings(tickers)
	return tickers
}

func (a *tickersStr) fetchTickerName(text string) []string {
	words := strings.Fields(text)
	words = words[1:]
	return words
}

//посмотрел отличие Marshal от MarshalIndent, подумал,
//что Indent лушче по читабельности - функциональной разницы нет
//обязательно на проверку с dataSave
func (a *tickersStr) writeInJson(tickersVault map[int]map[string]interface{}) {
	file, err := json.MarshalIndent(tickersVault, "", " ")
	if err != nil {
		log.Panic(err)
	}
	_ = ioutil.WriteFile("test.json", file, 0644)
	//defer wg.Done()

}

func (a *tickersStr) ReadFromJson(tickersVault map[int]map[string]interface{}) {
	file, err := ioutil.ReadFile("test.json")
	if err != nil {
		log.Panic(err)
	}
	json.Unmarshal(file, &tickersVault)
}

func (a *tickersStr) autoSave(tickersVault map[int]map[string]interface{}) {
	for {
		time.Sleep(1 * time.Minute)
		a.writeInJson(tickersVault)
	}
}

func (a *tickersStr) ReturnUserTickerList(reply string, tickers []string) string {
	if tickers == nil {
		reply = "Empty ticker list"
	}
	for _, v := range tickers {
		reply = reply + v + " "
	}
	return reply
}

func (a *tickersStr) dataSave(tickersVault map[int]map[string]interface{}) {
	var wg sync.WaitGroup
	wg.Add(1)
	go a.writeInJson(tickersVault)
	wg.Wait()
}

//в разработке
func AddCandleByTicker(reply string, UsID int, words []string, tickersVault map[int]map[string]interface{}) string {

	



	for _, s := range words {
		for _, v := range tickersSlice {
			if v != s { //если такого тикера не существует
				reply = "Unknown ticker name " + s + "."

			} else if v == s { //если такой тикер найден

				client := sdk.NewRestClient(*token)



				candles, err := client.Candles(nil, time.Now().AddDate(0, 0, -20), time.Now(), sdk.CandleInterval1Month, s)
				if err != nil {
					log.Fatalln(err)
				}


				tickersVault[UsID][s] = nil
				reply = "Ticker " + s + " saved."
				break
			}
		}
	}
	fmt.Println(tickersVault)
	return reply
}

func checkCandlesExistence (figiName string) []{

	// получение часовых свечей за последние 20 дней по инструменту s
	client := sdk.NewRestClient(*token)

	figiName :=

	candles, err := client.Candles(nil, time.Now().AddDate(0, 0, -20), time.Now(), sdk.CandleInterval1Month, figiName)
	if err != nil {
		log.Fatalln(err)
	}
	figis := client.Candles()
}


//___________________________________________________________________________________________________________________
//___________________________________________________________________________________________________________________
