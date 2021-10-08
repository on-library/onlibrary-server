package models

import (
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
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt 		time.Time	`json:"updated_at"`
}