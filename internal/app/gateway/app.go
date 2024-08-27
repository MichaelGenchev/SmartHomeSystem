package gatewayapp

import (
	"net/http"

	"github.com/MichaelGenchev/smart-home-system/internal/common/config"
	"github.com/MichaelGenchev/smart-home-system/internal/gateway/handlers"
	"github.com/MichaelGenchev/smart-home-system/internal/gateway/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	config     *config.Config
	userConn   *grpc.ClientConn
	deviceConn *grpc.ClientConn
}

func NewApp(cfg *config.Config) (*App, error) {
	userConn, err := grpc.NewClient(cfg.Server.UserGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	deviceConn, err := grpc.NewClient(cfg.Server.DeviceGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &App{
		config:     cfg,
		userConn:   userConn,
		deviceConn: deviceConn,
	}, nil
}

func (a *App) Run() error {

	mux := http.NewServeMux()

	//User routes
	userHandler := handlers.NewUserHandler(a.userConn)
	mux.Handle("/api/v1/users/", middleware.Chain(
		http.HandlerFunc(userHandler.ServeHTTP),
		middleware.RateLimit,
		middleware.Auth(a.config.Auth.JWTSecret),
	))

	//Device routes
	deviceHandler := handlers.NewDeviceHandler(a.deviceConn)
	mux.Handle("/api/v1/devices/", middleware.Chain(
		http.HandlerFunc(deviceHandler.ServeHTTP),
		middleware.RateLimit,
		middleware.Auth(a.config.Auth.JWTSecret),
	))

	return http.ListenAndServe(":"+a.config.Server.GatewayPort, mux)
}

func (a *App) Shutdown() {
	if a.userConn != nil {
		a.userConn.Close()
	}

	if a.deviceConn != nil {
		a.deviceConn.Close()
	}
}
