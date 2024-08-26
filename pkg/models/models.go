// pkg/models/models.go
package models

import (
	"time"
)

// Device represents a smart home device
type Device struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Type      string    `bson:"type" json:"type"`
	State     string    `bson:"state" json:"state"`
	UserID    string    `bson:"user_id" json:"user_id"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// User represents a user of the smart home system
type User struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"-"` // Password is not included in JSON responses
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// DeviceState represents the current state of a device
type DeviceState struct {
	DeviceID  string    `bson:"device_id" json:"device_id"`
	State     string    `bson:"state" json:"state"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

// DeviceType represents a type of smart home device
type DeviceType struct {
	ID          string `bson:"_id,omitempty" json:"id"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}

// Room represents a room in the smart home
type Room struct {
	ID        string   `bson:"_id,omitempty" json:"id"`
	Name      string   `bson:"name" json:"name"`
	DeviceIDs []string `bson:"device_ids" json:"device_ids"`
}

// Schedule represents a scheduled action for a device
type Schedule struct {
	ID       string    `bson:"_id,omitempty" json:"id"`
	DeviceID string    `bson:"device_id" json:"device_id"`
	Action   string    `bson:"action" json:"action"`
	Time     time.Time `bson:"time" json:"time"`
	Repeat   bool      `bson:"repeat" json:"repeat"`
}
