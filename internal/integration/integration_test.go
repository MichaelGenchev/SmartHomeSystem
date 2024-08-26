//go:build integration
// +build integration

package integration

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/MichaelGenchev/smart-home-system/internal/device/repository"
	"github.com/MichaelGenchev/smart-home-system/internal/device/service"
	"github.com/MichaelGenchev/smart-home-system/pkg/models"
	"github.com/MichaelGenchev/smart-home-system/pkg/proto"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IntegrationTestSuite struct {
	suite.Suite
	sqlDB       *sql.DB
	mongoClient *mongo.Client
	service     *service.DeviceService
}

func (suite *IntegrationTestSuite) SetupSuite() {
	// Connect to PostgreSQL
	postgresURI := os.Getenv("POSTGRES_URI")
	if postgresURI == "" {
		suite.T().Fatal("POSTGRES_URI is not set")
	}
	sqlDB, err := sql.Open("postgres", postgresURI)
	if err != nil {
		suite.T().Fatalf("Could not connect to PostgreSQL: %v", err)
	}
	suite.sqlDB = sqlDB

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		suite.T().Fatal("MONGO_URI is not set")
	}
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		suite.T().Fatalf("Could not connect to MongoDB: %v", err)
	}
	suite.mongoClient = mongoClient

	// Initialize repositories
	mongoDB := os.Getenv("MONGO_DB")
	if mongoDB == "" {
		suite.T().Fatal("MONGO_DB is not set")
	}
	sqlRepo := repository.NewSQLRepository(sqlDB)
	mongoRepo := repository.NewMongoRepository(mongoClient.Database(mongoDB).Collection("devices"))
	combinedRepo := repository.NewCombinedRepository(sqlRepo, mongoRepo)

	// Initialize service
	suite.service = service.NewDeviceService(combinedRepo)

	// Wait for databases to be ready
	suite.waitForDatabases()
}

func (suite *IntegrationTestSuite) waitForDatabases() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			suite.T().Fatal("Timeout waiting for databases")
		default:
			if err := suite.sqlDB.PingContext(ctx); err == nil {
				if err := suite.mongoClient.Ping(ctx, nil); err == nil {
					return
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	if suite.sqlDB != nil {
		suite.sqlDB.Close()
	}
	if suite.mongoClient != nil {
		suite.mongoClient.Disconnect(context.Background())
	}
}

func (suite *IntegrationTestSuite) TestCreateAndGetDevice() {
	ctx := context.Background()

	// Create a device
	createReq := &proto.CreateDeviceRequest{
		Name:   "Test Device",
		Type:   "Sensor",
		UserId: "user123",
	}
	createResp, err := suite.service.CreateDevice(ctx, createReq)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), createResp.Id)

	// Get the created device
	getReq := &proto.GetDeviceRequest{
		Id: createResp.Id,
	}
	getResp, err := suite.service.GetDevice(ctx, getReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createResp.Id, getResp.Id)
	assert.Equal(suite.T(), createReq.Name, getResp.Name)
	assert.Equal(suite.T(), createReq.Type, getResp.Type)
	assert.Equal(suite.T(), createReq.UserId, getResp.UserId)
}

func (suite *IntegrationTestSuite) TestUpdateDeviceState() {
	ctx := context.Background()

	// Create a device
	createReq := &proto.CreateDeviceRequest{
		Name:   "Test Device for Update",
		Type:   "Actuator",
		UserId: "user456",
	}
	createResp, err := suite.service.CreateDevice(ctx, createReq)
	assert.NoError(suite.T(), err)

	// Update the device state
	updateReq := &proto.UpdateDeviceStateRequest{
		Id:    createResp.Id,
		State: "On",
	}
	updateResp, err := suite.service.UpdateDeviceState(ctx, updateReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "On", updateResp.State)

	// Get the updated device
	getReq := &proto.GetDeviceRequest{
		Id: createResp.Id,
	}
	getResp, err := suite.service.GetDevice(ctx, getReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "On", getResp.State)
}

func (suite *IntegrationTestSuite) TestListDevices() {
	ctx := context.Background()

	// Create multiple devices
	for i := 0; i < 5; i++ {
		createReq := &proto.CreateDeviceRequest{
			Name:   fmt.Sprintf("List Test Device %d", i),
			Type:   "Sensor",
			UserId: "user789",
		}
		_, err := suite.service.CreateDevice(ctx, createReq)
		assert.NoError(suite.T(), err)
	}

	// List devices
	listReq := &proto.ListDevicesRequest{
		UserId:   "user789",
		Page:     1,
		PageSize: 10,
	}
	listResp, err := suite.service.ListDevices(ctx, listReq)
	assert.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), len(listResp.Devices), 5)
	assert.GreaterOrEqual(suite.T(), int(listResp.Total), 5)
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}