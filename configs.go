package main

import (
	"encoding/json"
	"io"
	"os"

	book_bot_database "github.com/RedBuld/book_bot_database"
	book_bot_rmq "github.com/RedBuld/book_bot_rmq"
)

func (DC *DownloadCenter) parseRMQConfig(filepath string) *book_bot_rmq.RMQ_Params {
	config := book_bot_rmq.RMQ_Params{}

	jsonFile, _ := os.Open(filepath)
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &config)

	return &config
}

func (DC *DownloadCenter) parseDBConfig(filepath string) *book_bot_database.DB_Params {
	config := book_bot_database.DB_Params{}

	jsonFile, _ := os.Open(filepath)
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &config)

	return &config
}

func (DC *DownloadCenter) parseQueueConfig(filepath string) *QueueConfig {
	config := QueueConfig{}

	jsonFile, _ := os.Open(filepath)
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &config)

	return &config
}
