package handler

import (
	userpb "CarStore/UserService/api/pb/user"
	"CarStore/UserService/internal/usecase"
	"context"
	"log"
)

type AuthHandler struct {
	userpb.UnimplementedUserServiceServer
	uc *usecase.UserUsecase
}

func NewAuthHandler(uc *usecase.UserUsecase) userpb.UserServiceServer {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) RegisterUser(ctx context.Context, req *userpb.RegisterUserRequest) (*userpb.AuthResponse, error) {
	log.Printf("RegisterUser request: %+v", req)
	token, err := h.uc.Register(ctx, req.Email, req.Username, req.Password, "user")
	if err != nil {
		return &userpb.AuthResponse{Status: err.Error()}, nil
	}
	return &userpb.AuthResponse{Token: token, Status: "OK"}, nil
}

func (h *AuthHandler) LoginUser(ctx context.Context, req *userpb.LoginUserRequest) (*userpb.AuthResponse, error) {
	log.Printf("LoginUser request: %+v", req)
	// pick identifier
	ident := req.Identifier
	token, err := h.uc.Login(ctx, ident, req.Password)
	if err != nil {
		return &userpb.AuthResponse{Status: err.Error()}, nil
	}
	return &userpb.AuthResponse{Token: token, Status: "OK"}, nil
}
