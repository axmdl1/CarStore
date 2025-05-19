package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"CarStore/UserService/pkg/jwt"
)

type contextKey string

const (
	ContextKeyUserID   contextKey = "userID"
	ContextKeyUserRole contextKey = "userRole"
)

// methodACL defines the required minimum role for each RPC.
// Roles: "anon", "user", "admin"; admin > user > anon
var methodACL = map[string]string{
	// UserService
	"/user.UserService/RegisterUser":         "anon",
	"/user.UserService/LoginUser":            "anon",
	"/user.UserService/SendVerificationCode": "anon",
	"/user.UserService/ConfirmEmail":         "anon",
	"/user.UserService/GetProfile":           "user",
	"/user.UserService/ListUsers":            "admin",
	"/user.UserService/ChangeUserRole":       "admin",

	// CarService
	"/car.CarService/ListCars":      "anon",
	"/car.CarService/GetCar":        "anon",
	"/car.CarService/CreateCar":     "admin",
	"/car.CarService/UpdateCar":     "admin",
	"/car.CarService/DeleteCar":     "admin",
	"/car.CarService/DecreaseStock": "user",

	// OrderService
	"/order.OrderService/CreateOrder": "user",
	"/order.OrderService/GetOrder":    "user",
	"/order.OrderService/ListOrders":  "admin",
	"/order.OrderService/UpdateOrder": "admin",
	"/order.OrderService/DeleteOrder": "admin",
}

// UnaryAuthInterceptor returns a gRPC interceptor enforcing JWT auth and role-based access.
func UnaryAuthInterceptor(jwtSvc jwt.JWTService) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		required, ok := methodACL[info.FullMethod]
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "method not allowed")
		}
		if required == "anon" {
			return handler(ctx, req)
		}
		md, _ := metadata.FromIncomingContext(ctx)
		auth := ""
		if vals := md.Get("authorization"); len(vals) > 0 {
			auth = vals[0]
		}
		if !strings.HasPrefix(auth, "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "missing or invalid authorization header")
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := jwtSvc.ValidateToken(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		role := claims.Role
		switch required {
		case "user":
			if role != "user" && role != "admin" {
				return nil, status.Error(codes.PermissionDenied, "user role required")
			}
		case "admin":
			if role != "admin" {
				return nil, status.Error(codes.PermissionDenied, "admin role required")
			}
		}
		// inject into context
		ctx = context.WithValue(ctx, ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyUserRole, role)
		return handler(ctx, req)
	}
}

// FromContext retrieves the userID and role from context.
func FromContext(ctx context.Context) (userID, role string) {
	if v := ctx.Value(ContextKeyUserID); v != nil {
		userID = v.(string)
	}
	if v := ctx.Value(ContextKeyUserRole); v != nil {
		role = v.(string)
	}
	return
}
