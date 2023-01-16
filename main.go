package main

import (
	"fmt"
	"log"
	"os"
	"time"

	book_bot_database "github.com/RedBuld/book_bot_database"
	book_bot_rmq "github.com/RedBuld/book_bot_rmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

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

	queue := DC.initQueue()
	defer DC.stopQueue()
	DC.queue = queue

	forever := make(chan bool)
	go DC.startQueue()
	// time.Sleep(5 * time.Second)
	<-forever
}

func (DC *DownloadCenter) initRMQ() *book_bot_rmq.RMQ_Session {
	config := DC.parseRMQConfig("configs/rabbitmq.json")
	config.Consumer = DC.onMessage
	rmq := book_bot_rmq.NewRMQ(config)

	return rmq
}

func (DC *DownloadCenter) initDB() *book_bot_database.DB_Session {
	config := DC.parseDBConfig("configs/database.json")
	db := book_bot_database.NewDB(config)

	return db
}

func (DC *DownloadCenter) initQueue() *Queue {
	config := DC.parseQueueConfig("configs/queue.json")
	queue := DC.newQueue(config)

	return queue
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
