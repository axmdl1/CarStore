package handler

import (
	userpb "CarStore/UserService/api/pb/user"
	"CarStore/UserService/internal/usecase"
	"CarStore/UserService/pkg/auth"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is required")
	}

	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "Username is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Password is required")
	}

	if err := h.uc.Register(ctx, req.Email, req.Username, req.Password, "user"); err != nil {
		return &userpb.AuthResponse{Status: err.Error()}, nil
	}
	return &userpb.AuthResponse{Token: "", Status: "code_sent"}, nil
}

func (h *AuthHandler) LoginUser(ctx context.Context, req *userpb.LoginUserRequest) (*userpb.AuthResponse, error) {
	log.Printf("LoginUser request: %+v", req.Identifier) //password hidden
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
	uid, _ := auth.FromContext(ctx)
	u, err := h.uc.Profile(ctx, uid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
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

func (h *AuthHandler) SendVerificationCode(ctx context.Context, req *userpb.SendCodeRequest) (*userpb.SendCodeResponse, error) {
	log.Printf("SendVerificationCode request: %+v", req)
	status, err := h.uc.SendVerificationCode(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return &userpb.SendCodeResponse{Status: status}, nil
}

func (h *AuthHandler) ConfirmEmail(ctx context.Context, req *userpb.ConfirmEmailRequest) (*userpb.ConfirmEmailResponse, error) {
	log.Printf("ConfirmEmail request: %+v", req)
	token, err := h.uc.ConfirmEmail(ctx, req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &userpb.ConfirmEmailResponse{Token: token, Status: "verified"}, nil
}
