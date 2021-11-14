package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Rent struct {
	ID				uuid.UUID	`gorm:"primaryKey" json:"id"`
	UserID			uuid.UUID	`gorm:"size:191" json:"user_id"`
	BookID			uuid.UUID	`gorm:"size:191" json:"book_id"`
	RentStatus		int			`json:"rent_status"`
	Fine			int			`json:"fine"`
	RentDescription	string		`json:"rent_description"`
	ExtendReason	string		`json:"extend_reason"`
	IsExtendConfirm	int			`json:"is_extend_confirm"`
	RentAt			time.Time	`json:"rent_at"`
	EndAt			time.Time	`json:"end_at"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt 		time.Time	`json:"updated_at"`
}