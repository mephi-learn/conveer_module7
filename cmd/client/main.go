package main

import (
	"context"
	dz7v1 "dz7/proto/gen"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var registerPort, secretPort, serverPort string

func main() {
	args := os.Args
	if len(args) < 3 {
		usage(filepath.Base(args[0]))
		return
	}

	registerPort = os.Getenv("REGISTER_PORT")
	if num, err := strconv.Atoi(registerPort); err != nil || num > 65535 || num < 1024 {
		registerPort = "53000"
	}

	secretPort = os.Getenv("SECRET_PORT")
	if num, err := strconv.Atoi(secretPort); err != nil || num > 65535 || num < 1024 {
		secretPort = "53001"
	}

	serverPort = os.Getenv("SERVER_PORT")
	if num, err := strconv.Atoi(serverPort); err != nil || num > 65535 || num < 1024 {
		serverPort = "55000"
	}

	// Получаем тип операции и её параметр
	operation := args[1]
	param := args[2]

	switch operation {
	case "register":
		fmt.Print(Register(param))
	case "secret":
		fmt.Print(Secret(param))
	case "server":
		mux := http.NewServeMux()
		mux.HandleFunc("/", getUsage)
		mux.HandleFunc("/register", getRegister)
		mux.HandleFunc("/secret", getSecret)

		err := http.ListenAndServe(":"+serverPort, mux)
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed\n")
		} else if err != nil {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}
	default:
		fmt.Print(usage(args[0]))
	}
}

func getUsage(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, usageHttp(fmt.Sprintf("%s%s", r.Host, r.URL.Path)))
}
func getRegister(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, Register(r.URL.Query().Get("secret")))
}
func getSecret(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, Secret(r.URL.Query().Get("id")))
}

func usage(programName string) string {
	result := fmt.Sprintf("Usage: %s command param\n", programName)
	result += fmt.Sprintf("%s register secret - register secret, return secret_id\n", programName)
	result += fmt.Sprintf("%s info secret_id - get secret by id\n", programName)
	return result
}

func usageHttp(host string) string {
	result := fmt.Sprintf("Usage: %scommand?data=param\n", host)
	result += fmt.Sprintf("%sregister?secret=value - register secret, return secret_id\n", host)
	result += fmt.Sprintf("%ssecret?id=secret_id - get secret by id\n", host)
	return result
}

func Register(param string) string {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial("server:"+registerPort, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	client := dz7v1.NewRegisterClient(conn)
	request := &dz7v1.RegisterRequest{
		Secret: param,
	}
	response, err := client.Register(context.Background(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	if response.GetError() != "" {
		return fmt.Sprintf("ERROR: register failure: %s\n", response.GetError())
	} else {
		return fmt.Sprintf("%s", response.GetSecretId())
	}
}

func Secret(param string) string {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial("server:"+secretPort, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	client := dz7v1.NewSecretClient(conn)
	request := &dz7v1.SecretRequest{
		SecretId: param,
	}
	response, err := client.Secret(context.Background(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	if response.GetError() != "" {
		return fmt.Sprintf("ERROR: receive failure: %s\n", response.GetError())
	} else {
		return fmt.Sprintf("%s\n", response.GetSecret())
	}
}
