package main

import (
	"time"

	"github.com/bigdataconcept/fintech/marketdataProvider/dataProvider"
	"github.com/bigdataconcept/fintech/marketdataProvider/threadPool"
	log "github.com/sirupsen/logrus"
	"github.com/bigdataconcept/fintech/marketdataProvider/util"
)

func main() {

	log.Println("Market Data Provider starting Get Live Stock Qoutes")
	pool := threadPool.NewThreadPool(1, 100000)
	symbols := dataProvider.NewSymbols()
	for {
		for _, v := range symbols.SymbolList {
			task := &StockQuoteTask{Symbol: v}
			pool.Execute(task)
		}
		time.Sleep(1000000)
	}

}

type StockQuoteTask struct {
	Symbol string
}

func (t *StockQuoteTask) Run() {
	marketDataProvider := &dataProvider.MarketDataProvider{}
	marketDataProvider.GetQouteMarketData(t.Symbol)
}



func getConfig() (util.AppConfig, error) {
	appConfig, err := util.LoadConfiguration("config.json")
	if err != nil {
		panic(err)
	}
	return appConfig, nil
}

