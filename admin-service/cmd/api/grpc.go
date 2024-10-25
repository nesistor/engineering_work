package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"admin-service/data"  
	"admin-service/admin" 
)

type AdminServer struct {
	admins.AdminServiceServer
	Models data.Models
}

func (s *AdminServer) ValidateAdmin(ctx context.Context, req *admins.ValidateAdminRequest) (*admins.ValidateAdminResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	// Pobranie administratora z bazy danych na podstawie emaila
	admin, err := s.Models.Admin.GetAdminByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return &admins.ValidateAdminResponse{
				IsValid: false,
				Message: "Admin not found",
			}, nil
		}
		return nil, err
	}

	// Weryfikacja hasła
	valid, err := s.Models.Admin.PasswordMatches(admin.ID, password)
	if err != nil || !valid {
		return &admins.ValidateAdminResponse{
			IsValid: false,
			Message: "Invalid password",
		}, nil
	}

	// Odpowiedź, że administrator jest poprawny
	return &admins.ValidateAdminResponse{
		IsValid: true,
		AdminId: admin.ID,
		Message: "Admin authenticated successfully",
	}, nil
}

// gRPCListen starts the gRPC server for the admin service
func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()
	admins.RegisterAdminServiceServer(s, &AdminServer{Models: app.Models})

	log.Printf("gRPC Server started on port %s", gRPCPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
