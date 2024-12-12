package main

import (
	dz7v1 "dz7/proto/gen"
	"fmt"
	"net"

	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	listener, err := net.Listen("tcp", ":53000")

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	dz7v1.RegisterHomeworkServer(grpcServer, &server{Secrets: map[string]string{}})
	grpcServer.Serve(listener)
}

type server struct {
	dz7v1.UnimplementedHomeworkServer
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
