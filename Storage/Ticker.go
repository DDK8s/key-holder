package Storage

import (

	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

func WriteInJson(tickersVault map[int]map[string]interface{}) {
	file, err := json.MarshalIndent(tickersVault, "", " ")
	if err != nil {
		log.Panic(err)
	}
	_ = ioutil.WriteFile("test.json", file, 0644)
	//defer wg.Done()

}

func ReadFromJson(tickersVault map[int]map[string]interface{}) {
	file, err := ioutil.ReadFile("test.json")
	if err != nil {
		log.Panic(err)
	}
	json.Unmarshal(file, &tickersVault)
}

func AutoSave(tickersVault map[int]map[string]interface{}) {
	for {
		time.Sleep(1 * time.Minute)
		WriteInJson(tickersVault)
	}
}

func CheckingMapInitialization(tickersVault map[int]map[string]interface{}, UsID int) {
	_, ok := tickersVault[UsID]
	if !ok {
		tickersVault[UsID] = make(map[string]interface{})
	}
}

func ReturnUserTickerList(reply string, tickers []string) string {
	if tickers == nil {
		reply = "Empty ticker list"
	}
	for _, v := range tickers {
		reply = reply + v + " "
	}
	return reply
}

func DataSave(tickersVault map[int]map[string]interface{}) {
	var wg sync.WaitGroup
	wg.Add(1)
	go WriteInJson(tickersVault)
	wg.Wait()
}

