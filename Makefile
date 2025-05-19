# Makefile for generating protobuf and grpc-gateway stubs for all services

# Directories
PROTO_DIR := proto
THIRD_PARTY := $(PROTO_DIR)/third_party/googleapis
USER_PB := UserService/api/pb
CAR_PB := CarService/api/pb
ORDER_PB := OrderService/api/pb

# Tools
PROTOC := protoc
GO_PLUGIN := plugins=grpc:
GW_PLUGIN := plugins=grpc-gateway:
GOOGLEAPIS := $(THIRD_PARTY)

# Include paths
INCLUDES := -I $(PROTO_DIR) -I $(GOOGLEAPIS)

# Proto files
USER_PROTO := $(PROTO_DIR)/user/user.proto
CAR_PROTO  := $(PROTO_DIR)/car/car.proto
ORDER_PROTO := $(PROTO_DIR)/order/order.proto

.PHONY: all user car order clean

all: user car order

# Generate stubs for UserService
user:
	@echo "Generating protobuf for UserService..."
	$(PROTOC) $(INCLUDES) \
	  --go_out=paths=source_relative:$(USER_PB) \
	  --go-grpc_out=paths=source_relative:$(USER_PB) \
	  --grpc-gateway_out=paths=source_relative:$(USER_PB) \
	  $(USER_PROTO)

# Generate stubs for CarService
car:
	@echo "Generating protobuf for CarService..."
	$(PROTOC) $(INCLUDES) \
	  --go_out=paths=source_relative:$(CAR_PB) \
	  --go-grpc_out=paths=source_relative:$(CAR_PB) \
	  --grpc-gateway_out=paths=source_relative:$(CAR_PB) \
	  $(CAR_PROTO)

# Generate stubs for OrderService
order:
	@echo "Generating protobuf for OrderService..."
	$(PROTOC) $(INCLUDES) \
	  --go_out=paths=source_relative:$(ORDER_PB) \
	  --go-grpc_out=paths=source_relative:$(ORDER_PB) \
	  --grpc-gateway_out=paths=source_relative:$(ORDER_PB) \
	  $(ORDER_PROTO)

# Clean generated files
clean:
	rm -rf $(USER_PB)/*.pb.go $(CAR_PB)/*.pb.go $(ORDER_PB)/*.pb.go
