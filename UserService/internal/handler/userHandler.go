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

func (h *AuthHandler) GetProfile(ctx context.Context, req *userpb.GetProfileRequest) (*userpb.ProfileResponse, error) {
	log.Printf("GetProfile request: %+v", req)
	u, err := h.uc.Profile(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return &userpb.ProfileResponse{User: &userpb.User{
		Id:       u.ID.String(),
		Email:    u.Email,
		Username: u.Username,
	}}, nil
}

func (h *AuthHandler) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	log.Printf("ListUsers request: %+v", req)
	us, err := h.uc.List(ctx)
	if err != nil {
		return nil, err
	}
	resp := &userpb.ListUsersResponse{}
	for _, u := range us {
		resp.Users = append(resp.Users, &userpb.User{
			Id:       u.ID.String(),
			Email:    u.Email,
			Username: u.Username,
			Password: u.Password,
		})
	}
	return resp, nil
}
