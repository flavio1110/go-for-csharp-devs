package main

import (
	"fmt"
	"log"
	"net"

	server "github.com/flavio1110/server_client/internal/MessageExchange"
	messageExchange "github.com/flavio1110/server_client/internal/server_client_pb"
	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", ":8082")

	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.ServerOption

	fmt.Println("Listening in ", ":8082")
	grpcServer := grpc.NewServer(opts...)
	messageExchange.RegisterMessageExchagerServiceServer(grpcServer, server.NewServer())
	grpcServer.Serve(lis)
}
