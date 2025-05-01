package service

import (
	"context"
	"log"

	"github.com/augustus281/beerus/proto/gen/go/beerus"
)

type UserService struct {
	beerus.UnimplementedUserServiceServer
}

func (u *UserService) AddUser(ctx context.Context, request *beerus.AddUserRequest) (*beerus.AddUserResponse, error) {
	log.Printf("Received request for adding a new user: %+v", request)
	return &beerus.AddUserResponse{
		Id: "user-id",
	}, nil
}
