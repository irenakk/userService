package service

import (
	"context"
	wallet "userService/proto/client"
)

type WalletGrpcService struct {
	client wallet.WalletServiceClient
}

func NewWalletService(client wallet.WalletServiceClient) *WalletGrpcService {
	return &WalletGrpcService{client: client}
}

func (r *WalletGrpcService) CreateWallet(userId int) (int, error) {
	response, err := r.client.CreateWallet(context.Background(), &wallet.WalletRequest{UserId: int64(userId)})
	if err != nil {
		return 0, err
	}

	return int(response.WalletId), nil
}

func (r *WalletGrpcService) CreateAccount(currency string, walletId int) (int, error) {
	response, err := r.client.CreateAccount(context.Background(), &wallet.AccountRequest{Currency: currency, WalletId: int64(walletId)})
	if err != nil {
		return 0, err
	}

	return int(response.AccountId), nil
}
