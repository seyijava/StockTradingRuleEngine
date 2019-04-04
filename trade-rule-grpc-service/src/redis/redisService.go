package redis

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type RedisTradeRuleService struct {
}

func NewPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

type TradeRule struct {
	AccountNumber string  `json:"accuntNumber"`
	Symbol        string  `json:"symbol"`
	SellLowPrice  float64 `json:"sellLowPrice"`
	SellHighPrice float64 `json:"sellHighPrice"`
	BuyLowPrice   float64 `json:"buyLowPrice"`
	BuyHighPrice  float64 `json:"buyHighPrice"`
	RuleType      string  `json:"ruleType"`
}

func (redisService *RedisTradeRuleService) SetTradeRule(c redis.Conn, tradeRule *TradeRule) error {

	const objectPrefix string = "TradeRule:"

	// serialize User object to JSON
	json, err := json.Marshal(tradeRule)
	if err != nil {
		return err
	}

	// SET object
	_, err = c.Do("SET", objectPrefix+tradeRule.AccountNumber, json)
	if err != nil {
		return err
	}

	return nil
}

func (redisService *RedisTradeRuleService) GetTradeRuleByAccountNumber(c redis.Conn, acctNumber string) (*TradeRule, error) {

	const objectPrefix string = "TradeRule:"
	s, err := redis.String(c.Do("GET", objectPrefix+acctNumber))
	if err == redis.ErrNil {
		fmt.Println("Rule does not exist")
	} else if err != nil {
		return nil, err
	}

	tradeRule := TradeRule{}
	err = json.Unmarshal([]byte(s), &tradeRule)

	fmt.Printf("%+v\n", tradeRule)

	return &tradeRule, err

}
