package main

import (
	dz7v1 "dz7/proto/gen"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	registerPort := os.Getenv("REGISTER_PORT")
	if num, err := strconv.Atoi(registerPort); err != nil || num > 65535 || num < 1024 {
		registerPort = "53000"
	}
	secretPort := os.Getenv("SECRET_PORT")
	if num, err := strconv.Atoi(secretPort); err != nil || num > 65535 || num < 1024 {
		secretPort = "53001"
	}

	registerListener, err := net.Listen("tcp", ":"+registerPort)
	secretListener, err := net.Listen("tcp", ":"+secretPort)

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	registerGrpcServer := grpc.NewServer(opts...)
	SecretGrpcServer := grpc.NewServer(opts...)
	srv := &server{Secrets: map[string]string{}}
	dz7v1.RegisterRegisterServer(registerGrpcServer, srv)
	go registerGrpcServer.Serve(registerListener)
	dz7v1.RegisterSecretServer(SecretGrpcServer, srv)
	go SecretGrpcServer.Serve(secretListener)

	stop := make(chan struct{})
	<-stop
}

type server struct {
	dz7v1.UnimplementedRegisterServer
	dz7v1.UnimplementedSecretServer
	Secrets map[string]string
}

func (s *server) Register(_ context.Context, request *dz7v1.RegisterRequest) (response *dz7v1.RegisterResponse, err error) {
	secret := request.GetSecret()
	secretId := ""
	for secretId == "" {
		newUUID, err := uuid.NewUUID()
		if err != nil {
			response = &dz7v1.RegisterResponse{
				Error: fmt.Sprintf("failed to create uuid: %s", err),
			}
			return response, err
		}
		if _, ok := s.Secrets[newUUID.String()]; !ok {
			secretId = newUUID.String()
		}
	}

	s.Secrets[secretId] = secret
	response = &dz7v1.RegisterResponse{
		SecretId: secretId,
	}

	return response, nil
}

func (s *server) Secret(_ context.Context, request *dz7v1.SecretRequest) (*dz7v1.SecretResponse, error) {
	secretId := request.GetSecretId()
	secret, ok := s.Secrets[secretId]
	if !ok {
		response := &dz7v1.SecretResponse{
			Error: fmt.Sprintf("secret not found: %s", secretId),
		}
		return response, nil
	}

	response := &dz7v1.SecretResponse{
		Secret: secret,
	}

	return response, nil
}
