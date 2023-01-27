package main

import (
	"log"
	"os"
	"time"
)

func (DC *DownloadCenter) newQueue(config *QueueConfig) *Queue {
	queue := Queue{
		config:        config,
		logger:        log.New(os.Stdout, "", log.LstdFlags),
		checkInterval: 1 * time.Second,
		done:          make(chan RunningQueueTask),
	}
	for group := range queue.config.Groups {
		queue.logger.Printf("Queue: QueueConfig_Group %+v\n", group)
	}
	return &queue
}

func (q *Queue) startQueue() {
	DC.logger.Println("Queue: Starting")
	for {
		DC.logger.Println("Queue: step")
		for task := range q.waiting {
			q.logger.Printf("Queue: Task %+v\n", task)
		}

		select {
		case <-q.done:
		case <-time.After(q.checkInterval):
		}
	}
}

func (q *Queue) stopQueue() {
	DC.logger.Println("Queue: Stopping")
	DC.done <- true
}
