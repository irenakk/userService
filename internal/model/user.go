package model

type User struct {
	ID               int    `json:"id"`
	Username         string `json:"username"`
	Password         string `json:"-"`
	Telegramnickname string `json:"telegram_nickname"`
	TelegramChatID   *int64 `json:"telegram_chat_id"`
}
