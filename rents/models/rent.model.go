package models

import (
	modelAuth "onlibrary/auth/models"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Rent struct {
	ID				uuid.UUID	`gorm:"primaryKey" json:"id"`
	UserID			modelAuth.Auth	`json:"user_id"`
	RentAt			time.Time	`json:"rent_at"`
	EndAt			time.Time	`json:"end_at"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt 		time.Time	`json:"updated_at"`
}