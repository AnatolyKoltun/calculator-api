//package main
//
//import (
//	"log"
//	"os"
//	"time"
//
//	"github.com/AnatolyKoltun/calculator-api/handlers"
//	pb "github.com/AnatolyKoltun/calculator-api/proto"
//	"github.com/gin-gonic/gin"
//	"github.com/nats-io/nats.go"
//	"google.golang.org/grpc"
//)

//import (
//	"github.com/gin-gonic/gin"
//
//	"github.com/AnatolyKoltun/calculator-api/handlers"
//)
//
//func setupAndRunServer() {
//	router := gin.Default()
//
//	router.POST("/calculate", handlers.CreateCalculation)
//	router.GET("/calculations", handlers.GetCalculations)
//
//	router.Run(":8080")
//}
//
//func main() {
//	setupAndRunServer()
//}

package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"

	"github.com/AnatolyKoltun/calculator-api/handlers"
	pb "github.com/AnatolyKoltun/calculator-api/proto"
)

func main() {
	// 1. Подключение к NATS
	natsURL := os.Getenv("NATS_URL")

	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal("Ошибка подключения к NATS:", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("Ошибка JetStream:", err)
	}

	// Создаем Stream (если не существует)
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "CALCULATIONS",
		Subjects: []string{"calculations.*"},
		Storage:  nats.FileStorage,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.Fatal("Ошибка создания Stream:", err)
	}

	// 2. Подключение gRPC клиента к storage
	storageURL := os.Getenv("STORAGE_SERVICE_URL")
	if storageURL == "" {
		storageURL = "localhost:50051"
	}

	grpcConn, err := grpc.Dial(storageURL, grpc.WithInsecure(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatal("Ошибка подключения к gRPC:", err)
	}
	defer grpcConn.Close()

	storageClient := pb.NewStorageServiceClient(grpcConn)

	// 3. Настройка Gin с передачей зависимостей в обработчики
	router := gin.Default()

	// POST /calculate — используем NATS publisher
	router.POST("/calculate", func(c *gin.Context) {
		// Вызываем обработчик с передачей JetStream
		handlers.CreateCalculationWithNATS(c, js)
	})

	// GET /calculations — используем gRPC client
	router.GET("/calculations", func(c *gin.Context) {
		// Вызываем обработчик с передачей gRPC клиента
		handlers.GetCalculationsWithGRPC(c, storageClient)
	})

	log.Println("API сервер запущен на :8080")
	router.Run(":8080")
}
