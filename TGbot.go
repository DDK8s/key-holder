package main

import (
	"awesomeProject/5/Storage"
	"flag"
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	//"awesomeProject/5/Buisnes"
	"context"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"
)


var token = flag.String("token", "", "NONE")

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

	var sls Storage.TickersInter = &Storage.TickersStr{}
	
	//var a tickersInt = &tickersStr{}
	var tickersVault = make(map[int]map[string]interface{})


	bot, err := tgbotapi.NewBotAPI("NONE")
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

	//поместил сюда L114-L126

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
			Storage.CheckingMapInitialization(tickersVault, UsID)
			words := fetchTickerName(text)
			reply = addTickers(reply, UsID, words, tickersVault)

		case "mytickers":
			tickers := sorting(tickersVault, UsID)

			reply = Storage.ReturnUserTickerList(reply,tickers)

		case "delete":
			words := fetchTickerName(text)
			reply = deleteTicker(reply, UsID, words, tickersVault)

		case "help":
			reply = help(reply)

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
		Storage.writeInJson(tickersVault)
		fmt.Println("signal received")
	}
}

//_____________________________Структуры_____________________________________________________________________________
/*
type tickersInt interface {
	help(string) string
	start(string) string
	fetchTickerName(string) []string
	ReturnUserTickerList(string, []string) string
	sorting(map[int]map[string]interface{}, int) []string
	addTickers(string, int, []string, map[int]map[string]interface{}) string
	deleteTicker(string, int, []string, map[int]map[string]interface{}) string
}*/

//_____________________________Методы________________________________________________________________________________


func start(reply string) string {
	reply = "Hello. I am your personal investment assistant. To find out what I can do, write \"/help\"."
	return reply
}

func help(reply string) string {
	reply = `◽Use the command "/addticker [ticker name]" to add a new ticker to your list of tickers.
◽Use the command "/delete [ticker name]" to remove the ticker from your list of tickers.
◽Use the command "/mytickers" to see a list of your tickers.`
	return reply
}

func addTickers(reply string, UsID int, words []string, tickersVault map[int]map[string]interface{}) string {
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

func deleteTicker(reply string, UsID int, words []string, tickersVault map[int]map[string]interface{}) string {
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

func sorting(tickersVault map[int]map[string]interface{}, UsID int) []string {
	tickers := make([]string, 0, len(tickersVault[UsID]))
	for v := range tickersVault[UsID] {
		tickers = append(tickers, v)
	}
	sort.Strings(tickers)
	return tickers
}

func fetchTickerName(text string) []string {
	words := strings.Fields(text)
	words = words[1:]
	return words
}

func AddCandleByTicker(reply string, UsID int, words []string, tickersVault map[int]map[string]interface{}, ctx context.Context) string {

	client := sdk.NewRestClient(*token)

	for _, s := range words {
		for _, v := range tickersSlice {//вместо tikcerSlice - []sdk.Instruments
			if v != s { //если такого тикера не существует
				reply = "Unknown ticker name " + s + "."

			} else if v == s { //если такой тикер найден

				instruments, err := client.InstrumentByTicker(ctx, s)
				if err != nil {
					log.Fatalln(err)
				}
				log.Printf("%+v\n", instruments)//нулевой или нет, и взять нужный и хранить в структуре отдельной

				instrument, err := client.InstrumentByFIGI(ctx, s)
				if err != nil {
					log.Fatalln(err)
				}
				log.Printf("%+v\n", instrument)


				/*	Что за cancel?
					ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()*/
				tickersVault[UsID][s] = nil
				reply = "Ticker " + s + " saved."
				break
			}
		}
	}
	fmt.Println(tickersVault)
	return reply
}
