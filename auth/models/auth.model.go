package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)


type Auth struct {
	ID				uuid.UUID	`gorm:"primaryKey" json:"id"`
	Username		string		`json:"username"`
	Password		string		`json:"password"`
	Role			string		`json:"role"`
	Name			string		`json:"name"`
	Email			string		`json:"email"`
	Address			string		`json:"address"`
	City			string		`json:"city"`
	Province		string		`json:"province"`		
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt 		time.Time	`json:"updated_at"`
}