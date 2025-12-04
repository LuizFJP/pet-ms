package main

import (
	"github.com/LuizFJP/pet-ms/application"
	"github.com/LuizFJP/pet-ms/infrastructure/persistence"
	server "github.com/LuizFJP/pet-ms/interfaces/grpc"
	pb "github.com/LuizFJP/pet-ms/proto"
	"log"
	"net"
	"os"

	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Config centraliza parâmetros de infra
type Config struct {
	DBDriver   string
	DBUser     string
	DBPassword string
	DBPort     string
	DBHost     string
	DBName     string
	GRPCAddr   string
}

// LoadConfig pode vir de env, flags, etc.
func LoadConfig() Config {
	return Config{
		DBDriver:   getEnv("DB_DRIVER", "postgres"),
		DBUser:     getEnv("DB_USER", "lgc_user"),
		DBPassword: getEnv("DB_PASSWORD", "lgc_teste_password"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBHost:     getEnv("DB_HOST", "pg_pet"),
		DBName:     getEnv("DB_NAME", "pet_db"),
		GRPCAddr:   getEnv("GRPC_ADDR", ":50051"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// bootstrapApp inicializa banco, automigrate e application layer.
// Isso aqui é facilmente mockável num teste.
func bootstrapApp(cfg Config) (*application.PetApplicationInterface, func(), error) {
	services, err := persistence.NewPetRepo(
		cfg.DBDriver,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBPort,
		cfg.DBHost,
		cfg.DBName,
	)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		services.Close()
	}

	err = services.Automigrate()
	if err != nil {
		return nil, nil, err
	}

	app := application.NewPetApplication(services.Pet)

	return &app, cleanup, nil
}

// newGRPCServer cria o servidor gRPC com interceptors, reflection e serviço registrado.
// Essa função é totalmente testável sem banco nem rede.
func newGRPCServer(app *application.PetApplicationInterface) *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpcprometheus.UnaryServerInterceptor),
		grpc.StreamInterceptor(grpcprometheus.StreamServerInterceptor),
	)

	// registra métricas padrão do gRPC
	grpcprometheus.Register(s)
	grpcprometheus.EnableHandlingTimeHistogram()

	// se quiser registrar métricas customizadas, pode usar prometheus.MustRegister(...)
	_ = prometheus.DefaultRegisterer // apenas pra lembrar que dá pra usar

	// reflection pro evans/grpcurl/grpcui
	reflection.Register(s)

	// registra seu serviço
	petServer := server.NewPetServer(*app)
	pb.RegisterPetServiceServer(s, petServer)

	return s
}

// startGRPCServer recebe um *grpc.Server e um endereço, faz listen e serve.
func startGRPCServer(s *grpc.Server, addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Printf("gRPC server listening at %v", lis.Addr())
	return s.Serve(lis)
}

func main() {
	cfg := LoadConfig()

	app, cleanup, err := bootstrapApp(cfg)
	if err != nil {
		log.Fatalf("failed to bootstrap application: %v", err)
	}
	defer cleanup()

	s := newGRPCServer(app)

	if err := startGRPCServer(s, cfg.GRPCAddr); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
}
