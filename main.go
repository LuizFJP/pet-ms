package main

import (
	"github.com/LuizFJP/pet-ms/application"
	"github.com/LuizFJP/pet-ms/infrastructure/persistence"
	server "github.com/LuizFJP/pet-ms/interfaces/grpc"
	pb "github.com/LuizFJP/pet-ms/proto"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	services, err := persistence.NewPetRepo("postgres", "lgc_user", "lgc_teste_password", "5432", "localhost", "pet_db")
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.Automigrate()

	app := application.NewPetApplication(services.Pet)

	s := grpc.NewServer()
	petServer := server.NewPetServer(app)
	pb.RegisterPetServiceServer(s, petServer)
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
