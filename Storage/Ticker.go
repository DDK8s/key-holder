package Storage

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

func (d *TickersStr) writeInJson(tickersVault map[int]map[string]interface{}) {
	file, err := json.MarshalIndent(tickersVault, "", " ")
	if err != nil {
		log.Panic(err)
	}
	_ = ioutil.WriteFile("test.json", file, 0644)
	//defer wg.Done()

}

func (d *TickersStr) readFromJson(tickersVault map[int]map[string]interface{}) {
	file, err := ioutil.ReadFile("test.json")
	if err != nil {
		log.Panic(err)
	}
	json.Unmarshal(file, &tickersVault)
}

func (d *TickersStr) autoSave(tickersVault map[int]map[string]interface{}) {
	for {
		time.Sleep(1 * time.Minute)
		d.writeInJson(tickersVault)
	}
}

func (d *TickersStr) checkingMapInitialization(tickersVault map[int]map[string]interface{}, UsID int) {
	_, ok := tickersVault[UsID]
	if !ok {
		tickersVault[UsID] = make(map[string]interface{})
	}
}

func (d *TickersStr) returnUserTickerList(reply string, tickers []string) string {
	if tickers == nil {
		reply = "Empty ticker list"
	}
	for _, v := range tickers {
		reply = reply + v + " "
	}
	return reply
}

func (d *TickersStr) dataSave(tickersVault map[int]map[string]interface{}) {
	var wg sync.WaitGroup
	wg.Add(1)
	go d.writeInJson(tickersVault)
	wg.Wait()

}

type TickersInter interface {
	writeInJson(map[int]map[string]interface{})
	readFromJson(map[int]map[string]interface{})
	autoSave(map[int]map[string]interface{})
	checkingMapInitialization(map[int]map[string]interface{}, int)
	returnUserTickerList(string, []string) string
	dataSave(map[int]map[string]interface{})
}

type TickersStr struct {
	tickersVault map[int]map[string]interface{}
}
