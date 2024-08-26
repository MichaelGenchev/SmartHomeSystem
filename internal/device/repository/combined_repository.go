package repository

import (
	"context"

	"github.com/MichaelGenchev/smart-home-system/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

type CombinedRepository struct {
	sqlRepo   *SQLRepository
	nosqlRepo *MongoRepository
}

func NewCombinedRepository(sqlRepo *SQLRepository, nosqlRepo *MongoRepository) *CombinedRepository {
	return &CombinedRepository{
		sqlRepo:   sqlRepo,
		nosqlRepo: nosqlRepo,
	}
}

func (r *CombinedRepository) CreateDevice(ctx context.Context, device *models.Device) (*models.Device, error) {
	// Write to SQL
	createdDevice, err := r.sqlRepo.CreateDevice(ctx, device)
	if err != nil {
		return nil, err
	}

	// Also write to NoSQL for read optimization
	_, err = r.nosqlRepo.collection.InsertOne(ctx, createdDevice)
	if err != nil {
		// Log error but don't fail the operation
		// In a production system, you might want to implement a retry mechanism or compensating transaction
		// TODO: Implement error logging
	}

	return createdDevice, nil
}

func (r *CombinedRepository) GetDevice(ctx context.Context, id string) (*models.Device, error) {
	// Read from NoSQL
	return r.nosqlRepo.GetDevice(ctx, id)
}

func (r *CombinedRepository) UpdateDeviceState(ctx context.Context, id string, state string) (*models.Device, error) {
	// Update in SQL
	updatedDevice, err := r.sqlRepo.UpdateDeviceState(ctx, id, state)
	if err != nil {
		return nil, err
	}

	// Also update in NoSQL
	_, err = r.nosqlRepo.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"state": state}},
	)
	if err != nil {
		// Log error but don't fail the operation
		// TODO: Implement error logging
	}

	return updatedDevice, nil
}

func (r *CombinedRepository) ListDevices(ctx context.Context, userID string, page, pageSize int) ([]*models.Device, int, error) {
	// Read from NoSQL
	return r.nosqlRepo.ListDevices(ctx, userID, page, pageSize)
}
