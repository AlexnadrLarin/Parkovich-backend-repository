package cmd

import (
	"log"
    "net"
	"net/http"

	"google.golang.org/grpc"

	"go-parkovich/microservices/analytics/internal"
	"go-parkovich/microservices/analytics/internal/api"
	"go-parkovich/microservices/analytics/internal/database"
	"go-parkovich/microservices/analytics/pkg/proto"
)

func Start() error {
    repo, err := database.NewUserEventsRepository()
    if err != nil {
        log.Printf("Ошибка подключения к ClickHouse: %v", err)
        return err
    }
    defer repo.Close() 

    lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Printf("Ошибка при прослушивании порта gRPC: %v", err)
		return err
	}

    grpcServer := grpc.NewServer()
	userEventServer := api.NewUserEventServer(repo, nil) 
	userevents.RegisterUserEventServiceServer(grpcServer, userEventServer)

	log.Println("gRPC сервер запущен на порту :50051")

	router := internal.SetupRouter(repo)

	go func() {
		log.Println("HTTP сервер микросервиса запущен на порту :8081")
		if err := http.ListenAndServe(":8081", router); err != nil {
			log.Printf("Ошибка при запуске HTTP сервера: %v", err)
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("Ошибка при запуске gRPC сервера: %v", err)
		return err
	}

	return nil
}