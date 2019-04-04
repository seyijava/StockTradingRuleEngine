package service

import (
	context "context"
	pb "proto"
	"redis"
)

type TradeRuleServer struct {
}

func (tradeRuleServer *TradeRuleServer) SubmitTradeRule(context context.Context, tradeRequest *pb.TradeRuleRequest) (*pb.TradeRuleResponse, error) {

	tradeRule := &redis.TradeRule{
		AccountNumber: tradeRequest.TradeRule.AcctNumber,
		Symbol:        tradeRequest.TradeRule.Symbol,
	}

	pool := redis.NewPool()
	conn := pool.Get()
	defer conn.Close()
	redisTradeRuleService := &redis.RedisTradeRuleService{}
	redisTradeRuleService.SetTradeRule(conn, tradeRule)
	return nil, nil
}
