package main

import (
	"time"
)

func (DC *DownloadCenter) newQueue(config *QueueConfig) *Queue {
	queue := Queue{
		config: config,
	}
	return &queue
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
