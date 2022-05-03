package model

type Setting struct {
	Port     int      `json:"port"`
	Interval string   `json:"interval"`
	Watch    []Watch  `json:"watch"`
	Logging  Logging  `json:"logging"`
	Telegram Telegram `json:"telegram"`
}

type Analytics struct {
	NowTime int64      `json:"now_time"`
	Results []Analytic `json:"results"`
}

// 메시지용 모델
type Analytic struct {
	TimeStamp int64  `json:"timestamp"`
	Label     string `json:"label"`
	Count     int    `json:"count"`
	Message   string `json:"message"`
}

// LogWatcher의 로그 설정
type Logging struct {
	Filename   string `json:"filename"`
	Level      string `json:"level"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

// 감시 대상
type Watch struct {
	Label    string   `json:"label"`
	Path     string   `json:"path"`
	File     string   `json:"file"`
	Level    []string `json:"level"`
	Contains []string `json:"contains"`
}

// 연동할 텔레그램 정보
type Telegram struct {
	BotToken string `json:"bot_token"`
	ChatId   int    `json:"chat_id"`
}
