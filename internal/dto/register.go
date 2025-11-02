package dto

type UserRegister struct {
	Username         string `json:"username" binding:"required"`
	Telegramnickname string `json:"telegram_nickname" binding:"required"`
	Password         string `json:"password" binding:"required,min=6"`
}
