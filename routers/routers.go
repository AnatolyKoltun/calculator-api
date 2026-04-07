package routers

import (
	"log"
	"os"

	"github.com/AnatolyKoltun/calculator-api/handlers"
	"github.com/AnatolyKoltun/calculator-api/proto"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

type CalculationServer struct {
	router        *gin.Engine
	natsPublisher nats.JetStreamContext
	grpcClient    proto.StorageServiceClient
}

func (s *CalculationServer) setupRoutes() {
	// POST /calculate — используем NATS publisher
	s.router.POST("/calculate", func(c *gin.Context) {
		handlers.CreateCalculationWithNATS(c, s.natsPublisher)
	})

	// GET /calculations — используем gRPC client
	s.router.GET("/calculations", func(c *gin.Context) {
		handlers.GetCalculationsWithGRPC(c, s.grpcClient)
	})
}

func (s *CalculationServer) Run() error {
	s.setupRoutes()

	port := os.Getenv("API_PORT")

	if port == "" {
		port = "8080"
	}

	err := s.router.Run(":" + port)

	if err == nil {
		log.Printf("API сервер запущен на %s", port)
	}

	return err
}

func NewCalculationServer(nats nats.JetStreamContext, grpc proto.StorageServiceClient) *CalculationServer {
	return &CalculationServer{
		router:        gin.Default(),
		natsPublisher: nats,
		grpcClient:    grpc,
	}
}
