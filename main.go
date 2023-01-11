package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	book_bot_database "github.com/RedBuld/book_bot_database"
	book_bot_rmq "github.com/RedBuld/book_bot_rmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type DownloadCenter struct {
	rmq    *book_bot_rmq.RMQ_Session
	logger *log.Logger
	db     *book_bot_database.DB_Session
	queue  *Queue
	done   chan bool
}

var DC *DownloadCenter

func main() {
	DC = &DownloadCenter{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		done:   make(chan bool),
	}

	rmq := DC.initRMQ()
	defer rmq.Close()
	DC.rmq = rmq

	db := DC.initDB()
	defer db.Close()
	DC.db = db

	config := parseConfigFromFile("config.json")
	queue := &Queue{
		config: config,
	}
	DC.queue = queue
	defer DC.stopQueue()
	go DC.startQueue()

	time.Sleep(5 * time.Second)
}

func (DC *DownloadCenter) initRMQ() *book_bot_rmq.RMQ_Session {
	params := &book_bot_rmq.RMQ_Params{
		Server: "amqp://guest:guest@localhost:5672/",
		Queue: &book_bot_rmq.RMQ_Params_Queue{
			Name:    "elib_fb2_downloads",
			Durable: true,
			// AutoAck: true,
		},
		Exchange: &book_bot_rmq.RMQ_Params_Exchange{
			Name:       "download_requests",
			Mode:       "topic",
			RoutingKey: "*",
			Durable:    true,
		},
		Prefetch: &book_bot_rmq.RMQ_Params_Prefetch{
			Count:  1,
			Size:   0,
			Global: false,
		},
		Consumer: DC.onMessage,
	}
	rmq := book_bot_rmq.NewRMQ(params)

	return rmq
}

func (DC *DownloadCenter) initDB() *book_bot_database.DB_Session {
	params := &book_bot_database.DB_Params{
		Server: "postgres://postgres:secret@localhost:5432/download-center",
	}
	db := book_bot_database.NewDB(params)

	return db
}

func parseConfigFromFile(filepath string) *QueueConfig {
	config := QueueConfig{}
	// var config map[string]interface{}

	jsonFile, _ := os.Open(filepath)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &config)

	return &config
}

func (DC *DownloadCenter) startQueue() {
	DC.logger.Println("Starting Queue")
	for {
		select {
		case <-DC.done:
			return
		case <-time.After(1 * time.Second):
		}
		DC.logger.Println("Queue step")
	}
}

func (DC *DownloadCenter) stopQueue() {
	DC.logger.Println("Stopping Queue")
	DC.done <- true
}

func (DC *DownloadCenter) onMessage(message amqp.Delivery) {
	fmt.Printf("[%s] Message [%s]: %s\n", time.Now(), message.RoutingKey, message.Body)
	message.Ack(false)
	// DC.SendStatus(message.RoutingKey)
}

// func (DC *DownloadCenter) SendStatus(RoutingKey string) {
// 	fmt.Println("Sending download status")

// 	message := &book_bot_rmq.RMQ_Message{
// 		Exchange:   "download_statuses",
// 		RoutingKey: RoutingKey,
// 		Mandatory:  false,
// 		Immediate:  false,
// 		Params: amqp.Publishing{
// 			DeliveryMode: amqp.Persistent,
// 			ContentType:  "text/plain",
// 			Body:         []byte("Download accepted"),
// 		},
// 	}
// 	err := DC.rmq.Push(message)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Sended download status")
// }
