package rpctransfer

import (
	"context"
	"userService/internal/model"
	"userService/proto/server"
)

func (h *Handlers) Find(ctx context.Context, in *server.FindRequest) (*server.FindResponse, error) {
	var user model.User
	user, err := h.Usecase.Find(in.Username)
	if err != nil {
		return &server.FindResponse{}, err
	}

	return &server.FindResponse{UserId: int64(user.ID), Username: user.Username}, nil
}
