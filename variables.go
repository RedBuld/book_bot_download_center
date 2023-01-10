package main

type Queue struct {
	config  *QueueConfig
	waiting []QueueTask
	active  []QueueTask
}

type QueueTask struct {
	//
}

type QueueConfig struct {
	Groups map[string]QueueConfig_Group `json:"groups"`
	Sites  map[string]QueueConfig_Site  `json:"sites"`
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

type DownloadRequest struct {
	userId    int64  `json:"user_id"`    // 123456789
	botId     string `json:"bot_id"`     // c1
	chatId    int64  `json:"chat_id"`    // 123456789
	messageId int64  `json:"message_id"` // 1200
	site      string `json:"site"`       // tl.rulate.ru
	url       string `json:"url"`        // https://tl.rulate.ru/book/71637
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
	Id    string // BOT_ID:CHAT_ID:MESSAGE_ID
	Text  string // Загружаю главу ...
	Files []string
}
