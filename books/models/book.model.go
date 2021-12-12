package models

import (
	categoryModel "onlibrary/category/models"
	genreModel "onlibrary/genre/models"
	reviewModel "onlibrary/reviews/models"
	"time"

	uuid "github.com/satori/go.uuid"
)


type Book struct {
	BookId			uuid.UUID					`gorm:"primaryKey" json:"id"`
	JudulBuku		string						`json:"judul_buku"`
	TahunTerbit		time.Time					`json:"tahun_terbit"`
	Penulis			string						`json:"penulis"`
	Penerbit		string						`json:"penerbit"`
	Stok			int							`json:"stok"`
	StokAwal		int							`json:"stok_awal"`
	Photo			string						`json:"photo"`
	DeskripsiBuku	string						`json:"deskripsi_buku"`
	ImgUrl			string						`json:"img_url"`
	Genres			[]genreModel.Genre			`gorm:"foreignKey:GenreBookID" json:"genres"`
	Reviews			[]reviewModel.Review	 	`gorm:"foreignKey:BookRefer" json:"reviews"`
	BookCategoryID	uuid.UUID					`gorm:"size:191" json:"-"`
	Category		categoryModel.Category		`gorm:"foreignKey:BookCategoryID" json:"category"`
	CreatedAt		time.Time					`json:"created_at"`
	UpdatedAt 		time.Time					`json:"updated_at"`
}