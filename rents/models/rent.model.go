package models

import (
	bookModel "onlibrary/books/models"
	"time"

	uuid "github.com/satori/go.uuid"
)


type Rent struct {
	PinjamID					uuid.UUID		`gorm:"primaryKey" json:"pinjam_id"`
	TanggalPinjam				time.Time		`json:"tanggal_pinjam"`
	TanggalPengembalian			time.Time		`json:"tanggal_pengembalian"`
	TanggalPengembalianFinish	*time.Time		`json:"tanggal_pengembalian_finish"`
	StatusPinjam				int				`json:"status_pinjam"`
	Denda						int				`json:"denda"`
	DeskripsiPeminjaman			string			`json:"deskripsi_peminjaman"`
	IsExtendConfirm				int				`json:"is_extend_confirm"`
	AlasanPerpanjangan			string			`json:"alasan_perpanjangan"`
	AuthID						uuid.UUID		`gorm:"size:191" json:"user_id"`
	BookRentID					uuid.UUID		`gorm:"size:191" json:"book_rent_id"`
	Book						bookModel.Book	`gorm:"foreignKey:BookRentID" json:"book"`
	CreatedAt					time.Time		`json:"created_at"`
	UpdatedAt 					time.Time		`json:"updated_at"`
}