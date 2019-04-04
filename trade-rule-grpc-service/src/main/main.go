package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "proto"
	"service"
)

const (
	port = ":7005"
)

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTradeRuleRpcServiceServer(s, &service.TradeRuleServer{})
	err = s.Serve(lis)
	if err != nil {
		fmt.Println("Failed")
		log.Fatalf("failed to serve: %v", err)
	} else {
		log.Println("RuleServer start Listener")
	}

}
