package dataProvider

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bigdataconcept/fintech/marketdataProvider/util"
	"github.com/optiopay/kafka"
	"github.com/optiopay/kafka/proto"
	"github.com/piquette/finance-go/quote"
	log "github.com/sirupsen/logrus"
)

var topic = "MarketDataTopic"
var partition = 1
var kafkaAddrs = make([]string, 1) //[]string{"192.168.2.24:9092"}

type MarketDataProvider struct {
	kafkaHost []string
	kafkaTopic string
	KafkaPartition int32
}

type StockQuote struct {
	Symbol                            string  `json:"symbol"`
	ShortName                         string  `json:"shortName"`
	QuoteType                         string  `json"quoteType"`
	RegularMarketChangePercent        float64 `json:"regularMarketChangePercent"`
	RegularMarketPreviousClose        float64 `json:"regularMarketPreviousClose"`
	RegularMarketPrice                float64 `json:"regularMarketPrice"`
	RegularMarketTime                 int64   `json:"regularMarketTime"`
	RegularMarketChange               float64 `json:"regularMarketChange"`
	RegularMarketOpen                 float64 `json:"regularMarketOpen"`
	RegularMarketDayHigh              float64 `json:"regularMarketDayHigh"`
	RegularMarketDayLow               float64 `json:"regularMarketDayLow"`
	RegularMarketVolume               int64   `json:"regularMarketVolume"`
	Bid                               float64 `json:"bid"`
	Ask                               float64 `json:"ask"`
	BidSize                           int64   `json:"bidSize"`
	AskSize                           int64   `json:"askSize"`
	PreMarketPrice                    float64 `json:"preMarketPrice"`
	PreMarketChangePercent            float64 `json:"preMarketChangePercent"`
	PreMarketTime                     int64   `json:"preMarketTime"`
	PostMarketPrice                   float64 `json:"preMarketTime"`
	PostMarketChange                  float64 `json:"postMarketChange"`
	PostMarketChangePercent           float64 `json:"postMarketChangePercent"`
	PostMarketTime                    int64   `json:"postMarketTime"`
	FiftyTwoWeekLowChange             float64 `json:"fiftyTwoWeekLowChange"`
	FiftyTwoWeekLowChangePercent      float64 `json:"fiftyTwoWeekLowChangePercent"`
	FiftyTwoWeekHighChange            float64 `json:"fiftyTwoWeekHighChange"`
	FiftyTwoWeekHighChangePercent     float64 `json:"fiftyTwoWeekHighChangePercent"`
	FiftyTwoWeekLow                   float64 `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHigh                  float64 `json:"fiftyTwoWeekHigh"`
	FiftyDayAverage                   float64 `json:"fiftyDayAverage"`
	FiftyDayAverageChange             float64 `json:"fiftyDayAverageChange"`
	FiftyDayAverageChangePercent      float64 `json:"fiftyDayAverageChangePercent"`
	TwoHundredDayAverage              float64 `json:"twoHundredDayAverage"`
	TwoHundredDayAverageChange        float64 `json:"twoHundredDayAverageChange"`
	TwoHundredDayAverageChangePercent float64 `json:"twoHundredDayAverageChangePercent"`
	AverageDailyVolume3Month          float64 `json:"averageDailyVolume3Month"`
	Currency                          string  `json:"currency"`
	Tradeable                         bool    `json:"tradeable"`
	ExchangeDataDelayedBy             int64   `json:"exchangeDataDelayedBy"`
	SourceInterval                    int64   `json:"sourceInterval"`
	ExchangeTimezoneName              string  `json:"exchangeTimezoneName"`
	ExchangeTimezoneShortName         string  `json:"exchangeTimezoneShortName"`
	Market                            string  `json:"market"`
	Exchange                          string  `json:"exchange"`
	QuoteTimeStamp                    string  `json:"quoteTimeStamp"`
}

func (marketDataProvider *MarketDataProvider) GetQuote(symbol string) (*StockQuote, error) {
	log.Println("Getting Stock Qoute for Stock Symbol " + symbol)
	q, err := quote.Get(symbol)
	if err != nil {
		panic(err)
	}
	payload, err := json.Marshal(q)
	bytes := []byte(string(payload))
	sq := &StockQuote{}
	err = json.Unmarshal(bytes, sq)
	if err != nil {
		panic(err)
	}

	layout := "2006-01-02T15:04:05"
	sq.QuoteTimeStamp = time.Now().Format(layout)
	return sq, err
}

func (marketDataProvider *MarketDataProvider) GetQouteMarketData(symbol string) error {
	//
	//
	quote, err := marketDataProvider.GetQuote(symbol)
	err = sendQuoteToKafka(quote)
	return err
}

func sendQuoteToKafka(quote *StockQuote) error {
	log.Println("Ingesting Stock Quote in Kafka Topic  " + topic)
	payload, err := json.Marshal(quote)
	quoteMessage := string(payload)

	fmt.Println(quoteMessage)
	conf := kafka.NewBrokerConf("test-client")
	conf.AllowTopicCreation = true
	broker, err := kafka.Dial(kafkaAddrs, conf)
	if err != nil {
		fmt.Printf("cannot connect to kafka cluster: %s", err)
	}

	produceMessageToKafka(broker, quoteMessage)

	return err
}

func produceMessageToKafka(broker kafka.Client, qoutePayload string) {
	conf := kafka.NewProducerConf()
	conf.RequiredAcks = proto.RequiredAcksLocal
	producer := broker.Producer(conf)

	msg := &proto.Message{Value: []byte(qoutePayload)}

	offval, err := producer.Produce(topic, 0, msg)

	if err != nil {
		log.Println(fmt.Sprintf("cannot produce message to [%s] [%s] [%s]", topic, partition, err))
	} else {
		log.Println()
		log.Println(fmt.Sprintf("offset [%d] Message Sent  [%s]", offval, qoutePayload))
	}
	defer broker.Close()

}

type Symbols struct {
	SymbolList map[int32]string
}

func NewSymbols() *Symbols {

	s := new(Symbols)
	s.SymbolList = map[int32]string{
		1:  "ORCL",
		2:  "GOOGL",
		3:  "AAPL",
		4:  "FB",
		5:  "AMZN",
		6:  "MSFT",
		7:  "IBM",
		8:  "TWTR",
		9:  "HPQ",
		10: "DELL",
	}
	return s
}

func init() {
	appConfig, err := util.LoadConfiguration("config.json")
	if err != nil {
		panic(err)
	}
	kafkaAddrs[0] = appConfig.KafkaConfig.BrokerHost
	topic = appConfig.KafkaConfig.Topic

}
