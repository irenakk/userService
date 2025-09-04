package model

type Account struct {
	ID       int    `json:"id"`
	Balance  int    `json:"balance"`
	Currency string `json:"currency"`
	WalletId int    `json:"wallet_id"`
}
