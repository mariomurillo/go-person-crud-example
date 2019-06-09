package grpc

import (
	"context"
	"google.golang.org/grpc"
	v1 "grpc-persona-crud-example/pkg/api/v1"
	"log"
	"net"
	"os"
	"os/signal"
)

func RunServer(ctx context.Context, v1Api v1.PersonaServiceServer, port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	v1.RegisterPersonaServiceServer(server, v1Api)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("shutting down gRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()
	log.Println("starting gRPC server...")
	return server.Serve(listener)
}