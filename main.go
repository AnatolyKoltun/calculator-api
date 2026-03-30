package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	pb "github.com/AnatolyKoltun/calculator-api/proto"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"

	"github.com/AnatolyKoltun/calculator-api/handlers"
)

func setupAndRunServer() {
	router := gin.Default()

	router.POST("/calculate", handlers.CreateCalculation)
	router.GET("/calculations", handlers.GetCalculations)

	router.Run(":8080")
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}

func main() {
	setupAndRunServer()

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

	// 3. HTTP handlers
	// POST /calculate — асинхронная отправка через NATS
	http.HandleFunc("/calculate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CalculationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Публикуем в NATS
		data, _ := json.Marshal(req)
		_, err := js.Publish("calculations.create", data)
		if err != nil {
			http.Error(w, "Failed to publish", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"status": "processing"})
	})

	// GET /calculation/{id} — синхронный запрос через gRPC
	http.HandleFunc("/calculation/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Path[len("/calculation/"):]
		if id == "" {
			http.Error(w, "ID required", http.StatusBadRequest)
			return
		}

		// gRPC запрос в storage
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		resp, err := storageClient.GetCalculation(ctx, &pb.GetRequest{Id: id})
		if err != nil {
			http.Error(w, "Storage service error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Println("API сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
