package models

import (
	uuid "github.com/satori/go.uuid"
)

type Genre struct {
	GenreID			uuid.UUID	`gorm:"primaryKey" json:"genre_id"`
	Nama			string		`json:"nama"`
	GenreBookID		uuid.UUID	`gorm:"size:191" json:"book_book_id"`
}