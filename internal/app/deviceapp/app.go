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
	"time"

	"github.com/MichaelGenchev/smart-home-system/internal/common/config"
	"github.com/MichaelGenchev/smart-home-system/internal/device/repository"
	"github.com/MichaelGenchev/smart-home-system/internal/device/service"
	"github.com/MichaelGenchev/smart-home-system/pkg/proto"
	_ "github.com/lib/pq" // PostgreSQL driver
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type App struct {
	cfg           *config.Config
	sqlDB         *sql.DB
	mongoClient   *mongo.Client
	grpcServer    *grpc.Server
	deviceService *service.DeviceService
}

func NewApp(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Initialize(ctx context.Context) error {
	var err error

	// Initialize SQL database connection
	a.sqlDB, err = sql.Open("postgres", a.cfg.PostgresURI)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Test SQL connection
	if err = a.sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	// Initialize MongoDB client
	a.mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(a.cfg.MongoURI))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test MongoDB connection
	if err = a.mongoClient.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Initialize repositories
	sqlRepo := repository.NewSQLRepository(a.sqlDB)
	mongoRepo := repository.NewMongoRepository(a.mongoClient.Database(a.cfg.MongoDB).Collection("devices"))
	combinedRepo := repository.NewCombinedRepository(sqlRepo, mongoRepo)

	// Initialize service
	a.deviceService = service.NewDeviceService(combinedRepo)

	// Initialize gRPC server
	a.grpcServer = grpc.NewServer()
	proto.RegisterDeviceServiceServer(a.grpcServer, a.deviceService)

	return nil
}

func (a *App) Run() error {
	// Start gRPC server
	listener, err := net.Listen("tcp", a.cfg.GRPCAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	go func() {
		log.Printf("Starting gRPC server on %s", a.cfg.GRPCAddr)
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.mongoClient.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect MongoDB: %v", err)
	} else {
		log.Println("MongoDB connection closed")
	}

	if err := a.sqlDB.Close(); err != nil {
		log.Printf("Failed to close PostgreSQL connection: %v", err)
	} else {
		log.Println("PostgreSQL connection closed")
	}

	return nil
}
