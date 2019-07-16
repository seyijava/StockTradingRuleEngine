package redis

import (
	"encoding/json"
	"fmt"

	pb "github.com/bigdataconcept/fintech/trade-rule/proto"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type RedisTradeRuleService struct {
}

const (
	RedisTradeRuleKey string = "TradeRule"
)

func RedisClient() *redis.Client {
	client := redis.NewClient(
		&redis.Options{
			Addr:     "192.168.2.24:6379",
			Password: "",
			DB:       0,
		})
	return client
}

type TradeRule struct {
	AccountNumber string  `json:"accountNumber"`
	Symbol        string  `json:"symbol"`
	SellPrice     float64 `json:"sellPrice"`
	BuyPrice      float64 `json:"buyPrice"`
	SellQuantity  int64   `json:"sellQuantity"`
	BuyQuantity   int64   `json:"buyQuantity"`
	RuleType      int32   `json:"ruleType"`
}

func (redisService *RedisTradeRuleService) SubmitTradeRule(redisClient *redis.Client, tradeRule *TradeRule) (int64, error) {
	log.Println("Submited Trade Rule from Client")
	json, err := json.Marshal(tradeRule)
	if err != nil {
		return 0, err
	}
	id, err := redisClient.LPush(RedisTradeRuleKey, json).Result()

	if err != nil {
		fmt.Println(err)
	} else if err != nil {
		return 0, err
	}
	return id, nil
}

func (redisService *RedisTradeRuleService) GetAllTradeRules(redisClient *redis.Client) (*pb.TradeRulesResponse, error) {
	log.Println("Retrive All Trade Rules from Redis")
	rules, err := redisClient.LRange(RedisTradeRuleKey, 0, -1).Result()
	log.Println(len(rules))
	if err != nil {
		log.Println("Rule does not exist")
	} else if err != nil {
		return nil, err
	}
	tradeRules := make([]*pb.TradeRule, len(rules))
	for index, rule := range rules {
		tradeRule := &pb.TradeRule{}
		err = json.Unmarshal([]byte(rule), tradeRule)
		tradeRules[index] = tradeRule

	}
	tradeResponse := &pb.TradeRulesResponse{TradeRules: tradeRules}
	log.Println(len(tradeResponse.TradeRules))
	return tradeResponse, nil

}
