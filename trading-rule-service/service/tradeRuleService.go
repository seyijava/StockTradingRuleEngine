package service

import (
	context "context"
	"strconv"

	pb "github.com/bigdataconcept/fintech/trade-rule/proto"
	"github.com/bigdataconcept/fintech/trade-rule/redis"
	log "github.com/sirupsen/logrus"
)

type TradeRuleServer struct {
}

func (tradeRuleServer *TradeRuleServer) SubmitTradeRule(context context.Context, tradeRequest *pb.TradeRuleRequest) (*pb.TradeRuleResponse, error) {
	log.Println("Trade Rules Submitted and Store in Redis")
	tradeRule := &redis.TradeRule{
		AccountNumber: tradeRequest.TradeRule.AcctNumber,
		Symbol:        tradeRequest.TradeRule.Symbol,
		SellPrice:     tradeRequest.TradeRule.SellPrice,
		BuyPrice:      tradeRequest.TradeRule.BuyPrice,
		RuleType:      tradeRequest.TradeRule.TradeType,
		SellQuantity:  tradeRequest.TradeRule.SellQuantity,
		BuyQuantity:   tradeRequest.TradeRule.BuyQuantity,
	}

	redisClient := redis.RedisClient()
	redisTradeRuleService := &redis.RedisTradeRuleService{}
	id, err := redisTradeRuleService.SubmitTradeRule(redisClient, tradeRule)
	if err != nil {
		return nil, err
	}
	ruleResponse := &pb.TradeRuleResponse{RuleNumber: strconv.FormatInt(id, 10)}
	return ruleResponse, err
}

func (tradeRuleServer *TradeRuleServer) GetTradeRules(ctx context.Context, tradeRequest *pb.TradeRulesRequest) (*pb.TradeRulesResponse, error) {
	log.Println("Get Trade Rules from Redis Data Store")
	redisClient := redis.RedisClient()
	redisTradeRuleService := new(redis.RedisTradeRuleService)
	rules, err := redisTradeRuleService.GetAllTradeRules(redisClient)
	if err != nil {
		log.Println("Rule does not exist")
	} else if err != nil {
		return nil, err
	}
	return rules, nil
}
