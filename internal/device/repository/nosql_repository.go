package repository

import (
	"context"

	"github.com/MichaelGenchev/smart-home-system/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(collection *mongo.Collection) *MongoRepository {
	return &MongoRepository{collection: collection}
}

func (r *MongoRepository) CreateDevice(ctx context.Context, device *models.Device) (*models.Device, error) {
	if device.ID == "" {
		device.ID = primitive.NewObjectID().Hex()
	}
	_, err := r.collection.InsertOne(ctx, device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (r *MongoRepository) GetDevice(ctx context.Context, id string) (*models.Device, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var device models.Device
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&device)
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *MongoRepository) UpdateDeviceState(ctx context.Context, id string, state string) (*models.Device, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{"$set": bson.M{"state": state}}
	var device models.Device
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&device)
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *MongoRepository) ListDevices(ctx context.Context, userID string, page, pageSize int) ([]*models.Device, int, error) {
	skip := (page - 1) * pageSize
	limit := int64(pageSize)

	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSkip(int64(skip)).SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var devices []*models.Device
	if err = cursor.All(ctx, &devices); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return devices, int(total), nil
}
