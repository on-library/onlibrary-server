package models

import (
	// categoryModel "onlibrary/category/models"
	genreModel "onlibrary/genre/models"
	reviewModel "onlibrary/reviews/models"
	"time"

	uuid "github.com/satori/go.uuid"
)


type Book struct {
	ID				uuid.UUID	`gorm:"primaryKey" json:"id"`
	BookId			uuid.UUID	`gorm:"primaryKey" json:"id"`
	JudulBuku		string		`json:"judul_buku"`
	TahunTerbit		time.Time	`json:"tahun_terbit"`
	Penulis			string		`json:"penulis"`
	Penerbit		string		`json:"penerbit"`
	Stok			int		`json:"stok"`
	Photo			string		`json:"photo"`
	DeskripsiBuku	string		`json:"deskripsi_buku"`
	// Title			string		`json:"title"`
	// Description		string		`json:"description"`
	// Author			string		`json:"author"`
	// Category		string		`json:"category"`
	// Publisher		string		`json:"publisher"`
	// Stock			uint		`json:"stock"`
	// Category		categoryModel.Category	`gorm:"foreignKey:BookID" json:"category"`
	Genres			[]genreModel.Genre		`gorm:"foreignKey:GenreBookID" json:"genres"`
	Reviews			[]reviewModel.Review	 `gorm:"foreignKey:BookRefer" json:"reviews"`
	BookCategoryID	uuid.UUID	`gorm:"size:191" json:"book_category_id"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt 		time.Time	`json:"updated_at"`
}