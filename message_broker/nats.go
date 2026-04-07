package message_broker

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

type HandleNats struct {
	nc *nats.Conn
}

func (n *HandleNats) CreateStreamNats() nats.JetStreamContext {
	var err error
	// 1. Подключение к NATS
	natsURL := os.Getenv("NATS_URL")

	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	n.nc, err = nats.Connect(natsURL)
	if err != nil {
		log.Fatal("Ошибка подключения к NATS:", err)
	}

	js, err := n.nc.JetStream()
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

	return js
}

func (n *HandleNats) Close() {
	n.nc.Close()
}
