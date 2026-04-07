package main

import (
	"log"

	"github.com/AnatolyKoltun/calculator-api/message_broker"
	"github.com/AnatolyKoltun/calculator-api/routers"
	"github.com/AnatolyKoltun/calculator-api/rpc"
)

func main() {
	var CalculationServer rpc.HandleGrpc
	var StreamNats message_broker.HandleNats

	defer CalculationServer.CloseConnectGrpcToStorage()
	defer StreamNats.Close()

	js := StreamNats.CreateStreamNats()
	storageClient := CalculationServer.ConnectGrpcToStorage()

	server := routers.NewCalculationServer(js, storageClient)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
