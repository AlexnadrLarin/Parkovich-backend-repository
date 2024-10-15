package main

import (
	"log"
	"net/http"

	"google.golang.org/grpc"

	"go-parkovich/microservices/analytics/cmd"
	"go-parkovich/microservices/analytics/pkg/proto"

)

func main() {
	go func() {
		if err := cmd.Start(); err != nil {
			log.Fatalf("Ошибка при запуске микросервиса: %v", err)
		}
	}()

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) 
	if err != nil {
		log.Fatalf("Ошибка подключения к микросервису по gRPC: %v", err)
	} 
	log.Println("Подключены к gRPC серверу на порту 50051")
	defer conn.Close()

	grpcClient := userevents.NewUserEventServiceClient(conn)

	_ = grpcClient

	log.Println("Главный сервер запущен на порту :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Ошибка запуска главного сервера: %v", err)
	}
}