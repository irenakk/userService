package usecase

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
	"userService/internal/dto"
	"userService/internal/model"
	"userService/internal/repository"
	"userService/internal/service"
	"userService/internal/utils"
)

type InterfaceUserUsecase interface {
	Create(ctx context.Context, user dto.UserRegister) (int, error)
	Find(username string) (model.User, error)
	Exists(username string) (bool, error)
	CheckPassword(loginPassword string, userPassword string) bool
	GenerateJWT(user model.User, tokenExpiration time.Duration, jwtSecret []byte) (string, error)
	LinkTelegramAccount(ctx context.Context, username string, chatID int64, tgnickname string) error
}

type UserUsecase struct {
	userRepository repository.InterfaceUserRepository
	walletService  service.WalletGrpcService
}

func NewUserUsecase(userRepository repository.InterfaceUserRepository, walletService service.WalletGrpcService) *UserUsecase {
	return &UserUsecase{userRepository, walletService}
}

func (usecase UserUsecase) Create(ctx context.Context, user dto.UserRegister) (int, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return 0, err
	}

	u := dto.UserRegister{
		Username: user.Username,
		Password: hashedPassword,
	}

	id, err := usecase.userRepository.Create(ctx, u)
	if err != nil {
		return 0, err
	}

	walletId, err := usecase.walletService.CreateWallet(id)
	if err != nil {
		usecase.walletService.DeleteWallet(walletId)
		usecase.userRepository.Delete(ctx, id)
		return 0, err
	}

	_, err = usecase.walletService.CreateAccount("USD", walletId)
	if err != nil {
		usecase.walletService.DeleteAccount("USD", walletId)
		usecase.walletService.DeleteWallet(walletId)
		usecase.userRepository.Delete(ctx, id)
		return 0, err
	}

	_, err = usecase.walletService.CreateAccount("EUR", walletId)
	if err != nil {
		usecase.walletService.DeleteAccount("EUR", walletId)
		usecase.walletService.DeleteAccount("USD", walletId)
		usecase.walletService.DeleteWallet(walletId)
		usecase.userRepository.Delete(ctx, id)
		return 0, err
	}

	_, err = usecase.walletService.CreateAccount("RUB", walletId)
	if err != nil {
		usecase.walletService.DeleteAccount("RUB", walletId)
		usecase.walletService.DeleteAccount("EUR", walletId)
		usecase.walletService.DeleteAccount("USD", walletId)
		usecase.walletService.DeleteWallet(walletId)
		usecase.userRepository.Delete(ctx, id)
		return 0, err
	}

	return id, nil
}

func (usecase UserUsecase) Find(username string) (model.User, error) {
	user, err := usecase.userRepository.Find(username)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (usecase UserUsecase) Exists(username string) (bool, error) {
	exists, err := usecase.userRepository.Exists(username)
	if err != nil {
		return true, err
	}
	return exists, nil
}

func (usecase UserUsecase) CheckPassword(loginPassword string, userPassword string) bool {
	return utils.CheckPasswordHash(loginPassword, userPassword)
}

func (usecase UserUsecase) GenerateJWT(user model.User, tokenExpiration time.Duration, jwtSecret []byte) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"username": user.Username,
		"iat":      now.Unix(),
		"exp":      now.Add(tokenExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, err
}

func (usecase *UserUsecase) LinkTelegramAccount(ctx context.Context, username string, chatID int64, tgnickname string) error {
	user, err := usecase.userRepository.Find(username)
	if err != nil {
		log.Println(err)
		return err
	}

	return usecase.userRepository.UpdateChatID(ctx, user.Username, chatID, tgnickname)
}
