package rpctransfer

import (
	"userService/internal/usecase"
	"userService/proto/server"
)

type Handlers struct {
	Usecase *usecase.UserUsecase
	server.UnimplementedUserServiceServer
}
