package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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

	// db := DC.initDB()
	// defer db.Close()
	// DC.db = db

	queue := DC.initQueue()
	defer queue.stopQueue()
	DC.queue = queue

	forever := make(chan bool)
	queue.startQueue()

	// DC.logger.Println("\nDC: QueueGroups BEFORE\n")
	// for group_name := range queue.groups {
	// 	DC.logger.Printf("DC: QueueGroups %s\n%+v\n", group_name, queue.groups[group_name])
	// }
	// time.Sleep(1 * time.Second)

	// DC.updateQueue()

	// DC.logger.Println("\nDC: QueueGroups AFTER\n")
	// for group_name := range queue.groups {
	// 	DC.logger.Printf("DC: QueueGroups %s\n%+v\n", group_name, queue.groups[group_name])
	// }

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

// func (DC *DownloadCenter) updateQueue() {
// 	config := DC.parseQueueConfig("configs/queue.json")
// 	qb := config.Groups["books"]
// 	qb.Simultaneously = 40
// 	config.Groups["books"] = qb
// 	DC.queue.updateConfig(config)
// }

func (DC *DownloadCenter) onMessage(message amqp.Delivery) {
	// fmt.Printf("[%s] Message [%s]: %s\n", time.Now(), message.RoutingKey, message.Body)
	var request DownloadRequest

	json.Unmarshal(message.Body, &request)
	fmt.Printf("Received download request:\n%+v\n", request)

	if DC.queue.addTask(request) {
		message.Ack(false)
		// go DC.SendSuccessStatus(request.BotId, request.ChatId, request.MessageId)
	} else {
		message.Nack(false, true)
		// go DC.SendFailStatus(request.BotId, request.ChatId, request.MessageId)
	}

}

func (DC *DownloadCenter) SendSuccessStatus(bot_id string, chat_id int64, message_id int64) {
	fmt.Println("Sending download status")

	status := DownloadStatus{
		BotId:     bot_id,
		ChatId:    chat_id,
		MessageId: message_id,
		Text:      "Download accepted",
		Files:     nil,
	}

	msg, err := json.Marshal(status)
	if err != nil {
		panic(err)
	}

	message := &book_bot_rmq.RMQ_Message{
		Exchange:   "download_statuses",
		RoutingKey: bot_id,
		Mandatory:  false,
		Immediate:  false,
		Params: amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         msg,
		},
	}
	err = DC.rmq.Push(message)
	if err != nil {
		panic(err)
	}

	fmt.Println("Sended download status")
}

func (DC *DownloadCenter) SendFailStatus(bot_id string, chat_id int64, message_id int64) {
	fmt.Println("Sending download status")

	status := DownloadStatus{
		BotId:     bot_id,
		ChatId:    chat_id,
		MessageId: message_id,
		Text:      "Download not accepted",
		Files:     nil,
	}

	msg, err := json.Marshal(status)
	if err != nil {
		panic(err)
	}

	message := &book_bot_rmq.RMQ_Message{
		Exchange:   "download_statuses",
		RoutingKey: bot_id,
		Mandatory:  false,
		Immediate:  false,
		Params: amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         msg,
		},
	}
	err = DC.rmq.Push(message)
	if err != nil {
		panic(err)
	}

	fmt.Println("Sended download status")
}
