// internal/app/deviceapp/app.go
package deviceapp

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/MichaelGenchev/smart-home-system/internal/common/config"
	grpcServer "github.com/MichaelGenchev/smart-home-system/internal/user/grpc"
	"github.com/MichaelGenchev/smart-home-system/internal/user/repository"
	"github.com/MichaelGenchev/smart-home-system/internal/user/service"

	"github.com/MichaelGenchev/smart-home-system/pkg/proto"
	_ "github.com/lib/pq" // PostgreSQL driver
	"google.golang.org/grpc"
)

type App struct {
	cfg         *config.Config
	sqlDB       *sql.DB
	grpcServer  *grpc.Server
	userService *grpcServer.GRPCServer
}

func NewApp(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Initialize(ctx context.Context) error {
	var err error

	// Initialize SQL database connection
	a.sqlDB, err = sql.Open("postgres", a.cfg.Database.UserPostgresURI)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Test SQL connection
	if err = a.sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	// Initialize repositories
	sqlRepo := repository.NewPostgresRepository(a.sqlDB)

	// Initialize service
	userservice := service.NewUserService(sqlRepo)

	// Initialize grpc server
	a.userService = grpcServer.NewGRPCServer(userservice)

	// Initialize gRPC server
	a.grpcServer = grpc.NewServer()
	proto.RegisterUserServiceServer(a.grpcServer, a.userService)

	return nil
}

func (a *App) Run() error {
	// Start gRPC server
	listener, err := net.Listen("tcp", a.cfg.Server.UserGRPCAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	go func() {
		log.Printf("Starting gRPC server on %s", a.cfg.Server.UserGRPCAddr)
		if err := a.grpcServer.Serve(listener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	a.grpcServer.GracefulStop()
	log.Println("gRPC server stopped")

	if err := a.sqlDB.Close(); err != nil {
		log.Printf("Failed to close PostgreSQL connection: %v", err)
	} else {
		log.Println("PostgreSQL connection closed")
	}

	return nil
}
