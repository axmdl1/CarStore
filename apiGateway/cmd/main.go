package main

import (
	"context"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	carpetpb "CarStore/CarService/api/pb/car"
	orderpb "CarStore/OrderService/api/pb/order"
	userpb "CarStore/UserService/api/pb/user" // adjust to your module path
)

func run() error {
	_ = godotenv.Load()

	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// register each service
	if err := userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:50052", opts); err != nil {
		return err
	}
	if err := carpetpb.RegisterCarServiceHandlerFromEndpoint(ctx, mux, "localhost:50053", opts); err != nil {
		return err
	}
	if err := orderpb.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, "localhost:50054", opts); err != nil {
		return err
	}

	var port = os.Getenv("API_GATEWAY_PORT")
	log.Println("Server listening on :" + port)
	return http.ListenAndServe(":"+port, mux)
}

func main() {
	flag.Parse()
	log.Fatal(run())
}
