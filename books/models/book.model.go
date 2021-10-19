package models

import (
	"onlibrary/reviews/models"
	"time"

	uuid "github.com/satori/go.uuid"
)


type Book struct {
	ID				uuid.UUID	`gorm:"primaryKey" json:"id"`
	Title			string		`json:"title"`
	Description		string		`json:"description"`
	Author			string		`json:"author"`
	Category		string		`json:"category"`
	Publisher		string		`json:"publisher"`
	Stock			uint		`json:"stock"`
	Reviews			[]models.Review	 `gorm:"foreignKey:BookRefer" json:"reviews"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt 		time.Time	`json:"updated_at"`
}