package repository

import (
	"context"
	"database/sql"

	"github.com/MichaelGenchev/smart-home-system/pkg/models"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) CreateDevice(ctx context.Context, device *models.Device) (*models.Device, error) {
	query := `INSERT INTO devices (name, type, state, user_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, device.Name, device.Type, device.State, device.UserID).Scan(&device.ID)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (r *SQLRepository) UpdateDeviceState(ctx context.Context, id string, state string) (*models.Device, error) {
	query := `UPDATE devices SET state = $1 WHERE id = $2 RETURNING id, name, type, state, user_id`
	var device models.Device
	err := r.db.QueryRowContext(ctx, query, state, id).Scan(&device.ID, &device.Name, &device.Type, &device.State, &device.UserID)
	if err != nil {
		return nil, err
	}
	return &device, nil
}
