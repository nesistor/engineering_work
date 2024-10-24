package main

import (
	"context"
	"fmt"
	"log"
	"user-service/data"
	"user-service/users"
	"net"
	"database/sql"

	"google.golang.org/grpc"
	
)

type UserServer struct {
	users.UnimplementedUserServiceServer
	Models data.Models
}

func (s *UserServer) ValidateUser(ctx context.Context, req *users.ValidateUserRequest) (*users.ValidateUserResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	// Pobranie użytkownika z bazy danych na podstawie emaila
	user, err := s.Models.User.GetUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return &users.ValidateUserResponse{
				IsValid: false,
				Message: "User not found",
			}, nil
		}
		return nil, err
	}

	// Weryfikacja hasła
	valid, err := s.Models.User.PasswordMatches(user.ID, password)
	if err != nil || !valid {
		return &users.ValidateUserResponse{
			IsValid: false,
			Message: "Invalid password",
		}, nil
	}

	// Odpowiedź, że użytkownik jest poprawny
	return &users.ValidateUserResponse{
		IsValid: true,
		UserId:  user.ID,
		Message: "User authenticated successfully",
	}, nil
}

// gRPCListen starts the gRPC server for the user service
func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()
	users.RegisterUserServiceServer(s, &UserServer{Models: app.Models})

	log.Printf("gRPC Server started on port %s", gRPCPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

