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
	id         int
	payload    *DownloadRequest
	downloader *Downloader
}

type Downloader struct {
	payload  *DownloadRequest
	log_file string
}

type DownloadRequest struct {
	userId    int64  `json:"user_id"`    // 123456789
	botId     string `json:"bot_id"`     // c1
	chatId    int64  `json:"chat_id"`    // 123456789
	messageId int64  `json:"message_id"` // 1200
	site      string `json:"site"`       // tl.rulate.ru
	url       string `json:"url"`        // https://tl.rulate.ru/book/7path
	start     int    `json:"start"`      // 0
	end       int    `json:"end"`        // 100
	format    string `json:"format"`     // fb2 [fb2|epub|cbz|mobi|azw3]
	login     string `json:"login"`      // zhena@zhizni.net
	password  string `json:"password"`   // ihatemyself
	images    bool   `json:"images"`     // true
	cover     bool   `json:"cover"`      // false
	proxy     string `json:"proxy"`      // 0.0.0.0:3128
}

type DownloadStatus struct {
	botId     string   `json:"bot_id"`     // c1
	chatId    int64    `json:"chat_id"`    // 123456789
	messageId int64    `json:"message_id"` // 1200
	Text      string   `json:"string"`     // Загружаю главу ...
	Files     []string `json:"files"`
}
