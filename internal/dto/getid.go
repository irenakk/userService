package dto

type Getid struct {
	Username string `json:"username" binding:"required"`
}
