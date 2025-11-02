package dto

type Link struct {
	Username   string `json:"username" binding:"required"`
	ChatId     int    `json:"chat_id" binding:"required"`
	Tgnickname string `json:"tg_nickname" binding:"required"`
}
