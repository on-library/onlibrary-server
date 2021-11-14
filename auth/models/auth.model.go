package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)


type Auth struct {
	ID				uuid.UUID	`gorm:"primaryKey" json:"id"`
	Username		string		`json:"username"`
	Password		string		`json:"-"`
	Role			int			`json:"role"`
	Name			string		`json:"name"`
	Email			string		`json:"email"`
	Address			string		`json:"address"`
	City			string		`json:"city"`
	Province		string		`json:"province"`
	IsVerify		int			`json:"is_verify"`
	VerifyCode		int			`json:"verify_code"`					
	RentedBook		string		`json:"rented_book"`		
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt 		time.Time	`json:"updated_at"`
}