package rpc

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/AnatolyKoltun/calculator-api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type HandleGrpc struct {
	grpcConn *grpc.ClientConn
	cancel   context.CancelFunc
}

func (g *HandleGrpc) ConnectGrpcToStorage() pb.StorageServiceClient {
	var err error
	var ctx context.Context
	// 1. Подключение gRPC клиента к storage
	storageURL := os.Getenv("STORAGE_SERVICE_URL")

	if storageURL == "" {
		storageURL = "localhost:50051"
	}

	// grpc.NewClient - неблокирующий
	g.grpcConn, err = grpc.NewClient(storageURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("Ошибка создания gRPC клиента:", err)
	}

	// Connect() НЕ возвращает ошибку! Он просто инициирует подключение.
	// Нужно использовать WaitForStateChange или проверять состояние.
	g.grpcConn.Connect()

	// Дожидаемся готовности соединения с таймаутом
	ctx, g.cancel = context.WithTimeout(context.Background(), 10*time.Second)

	// Правильный способ дождаться подключения
	for {
		state := g.grpcConn.GetState()
		if state == connectivity.Ready {
			break
		}
		if !g.grpcConn.WaitForStateChange(ctx, state) {
			log.Fatal("Таймаут подключения к gRPC серверу")
		}
	}

	log.Println("gRPC подключение к storage установлено")

	return pb.NewStorageServiceClient(g.grpcConn)
}

func (g *HandleGrpc) CloseConnectGrpcToStorage() {
	g.grpcConn.Close()
	g.cancel()
}
