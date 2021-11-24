package models

import (
	uuid "github.com/satori/go.uuid"
	// bookModel "onlibrary/books/models"
)

type Category struct {
	CategoryID			uuid.UUID		`gorm:"primaryKey" json:"category_id"`
	Nama			string				`json:"nama"`
	// Books			[]bookModel.Book	`gorm:"foreignKey:BookCategoryID" json:"books"`		
// 
}