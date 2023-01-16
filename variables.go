package main

import (
	"log"

	book_bot_database "github.com/RedBuld/book_bot_database"
	book_bot_rmq "github.com/RedBuld/book_bot_rmq"
)

type QueueConfig struct {
	Groups     map[string]QueueConfig_Group `json:"groups"`
	Sites      map[string]QueueConfig_Site  `json:"sites"`
	ExecFolder string                       `json:"exec_folder"`
}

type QueueConfig_Group struct {
	PerUser        int `json:"per_user"`
	Simultaneously int `json:"simultaneously"`
}

type QueueConfig_Site struct {
	Active         bool     `json:"active"`
	Parameters     []string `json:"parameters"`
	Proxy          string   `json:"proxy"`
	Simultaneously int      `json:"simultaneously"`
	PerUser        int      `json:"per_user"`
	Group          string   `json:"group"`
	PauseByUser    int      `json:"pause_by_user"`
}

type DownloadCenter struct {
	rmq    *book_bot_rmq.RMQ_Session
	logger *log.Logger
	db     *book_bot_database.DB_Session
	queue  *Queue
	done   chan bool
}

type Queue struct {
	config  *QueueConfig
	waiting []QueueTask
	active  []QueueTask
}

type QueueTask struct {
	Id         int
	Payload    *DownloadRequest
	Downloader *Downloader
}

type Downloader struct {
	Payload *DownloadRequest
	Log     string
}

type DownloadRequest struct {
	UserId    int64  `json:"user_id"`    // 123456789
	BotId     string `json:"bot_id"`     // c1
	ChatId    int64  `json:"chat_id"`    // 123456789
	MessageId int64  `json:"message_id"` // 1200
	Site      string `json:"site"`       // tl.rulate.ru
	Url       string `json:"url"`        // https://tl.rulate.ru/book/7path
	Start     int    `json:"start"`      // 0
	End       int    `json:"end"`        // 100
	Format    string `json:"format"`     // fb2 [fb2|epub|cbz|mobi|azw3]
	Login     string `json:"login"`      // zhena@zhizni.net
	Password  string `json:"password"`   // ihatemyself
	Images    bool   `json:"images"`     // true
	Cover     bool   `json:"cover"`      // false
	Proxy     string `json:"proxy"`      // 0.0.0.0:3128
}

type DownloadStatus struct {
	BotId     string   `json:"bot_id"`     // c1
	ChatId    int64    `json:"chat_id"`    // 123456789
	MessageId int64    `json:"message_id"` // 1200
	Text      string   `json:"string"`     // Загружаю главу ...
	Files     []string `json:"files"`
}
