package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Review struct {
	ID				uuid.UUID	`gorm:"primaryKey" json:"id"`
	Comment			string		`json:"comment"`
	Rating			uint		`json:"rating"`
	BookRefer		uuid.UUID	`gorm:"size:191" json:"book_refer"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt 		time.Time	`json:"updated_at"`
}