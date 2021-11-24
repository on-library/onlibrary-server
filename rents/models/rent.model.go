package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Rent struct {
	ID							uuid.UUID	`gorm:"primaryKey" json:"id"`
	// PinjamID					uuid.UUID	`gorm:"primaryKey" json:"pinjam_id"`
	// TanggalPinjam				time.Time	`json:"tanggal_pinjam"`
	// TanggalPengembalian			time.Time	`json:"tanggal_pengembalian"`
	// TanggalPengembalianFinish	time.Time	`json:"tanggal_pengembalian_finish"`
	// StatusPinjam				int			`json:"status_pinjam"`
	// Denda						int			`json:"denda"`
	// DeskripsiPeminjaman			string		`json:"deskripsi_peminjaman"`
	// IsExtendConfirm				int			`json:"is_extend_confirm"`
	// AlasanPerpanjangan			string		`json:"alasan_perpanjangan"`
	

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