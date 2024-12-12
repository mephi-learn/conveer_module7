package main

import (
	"context"
	dz7v1 "dz7/proto/gen"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"os"
	"path/filepath"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	args := os.Args
	conn, err := grpc.Dial("127.0.0.1:53000", opts...)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	if len(args) < 3 {

		usage(filepath.Base(args[0]))
		return
	}

	// Получаем тип операции и её параметр
	operation := args[1]
	param := args[2]

	client := dz7v1.NewHomeworkClient(conn)
	switch operation {
	case "register":
		request := &dz7v1.RegisterRequest{
			Secret: param,
		}
		response, err := client.Register(context.Background(), request)
		if err != nil {
			grpclog.Fatalf("fail to dial: %v", err)
		}

		if response.GetError() != "" {
			fmt.Printf("register failure: %s\n", response.GetError())
		} else {
			fmt.Printf("register success, your secret_id = %s\n", response.GetSecretId())
		}

	case "secret":
		request := &dz7v1.SecretRequest{
			SecretId: param,
		}
		response, err := client.Secret(context.Background(), request)
		if err != nil {
			grpclog.Fatalf("fail to dial: %v", err)
		}
		if response.GetError() != "" {
			fmt.Printf("receive failure: %s\n", response.GetError())
		} else {
			fmt.Printf("receive success, your secret: %s\n", response.GetSecret())
		}

	default:
		usage(args[0])
	}
}

func usage(programName string) {
	fmt.Printf("Usage: %s command param\n", programName)
	fmt.Printf("%s register secret - register secret, return secret_id\n", programName)
	fmt.Printf("%s info secret_id - get secret by id\n", programName)
}
