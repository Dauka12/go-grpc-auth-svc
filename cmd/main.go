package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Dauka12/go-grpc-auth-svc/pkg/config"
	"github.com/Dauka12/go-grpc-auth-svc/pkg/db"
	"github.com/Dauka12/go-grpc-auth-svc/pkg/pb"
	"github.com/Dauka12/go-grpc-auth-svc/pkg/services"
	"github.com/Dauka12/go-grpc-auth-svc/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "authorization_microservice",
		ExpirationHours: 24 * 365,
	}

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Auth Svc on", c.Port)

	s := services.Server{
		H:   h,
		Jwt: jwt,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
