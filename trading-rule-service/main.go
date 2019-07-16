package main

import (
	"net"

	pb "github.com/bigdataconcept/fintech/trade-rule/proto"
	"github.com/bigdataconcept/fintech/trade-rule/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	port = ":7005"
)

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTradeRuleRpcServiceServer(s, &service.TradeRuleServer{})
	log.Println("Trade Rule Server Started. Listining on port" + port)
	err = s.Serve(lis)
	if err != nil {
		log.Println("Trade Rule Server failed to startup")
		log.Errorf("failed to serve: %v", err)

	}

}
