package models

import (
	"time"

	rentsModel "onlibrary/rents/models"
	reviewModel "onlibrary/reviews/models"

	uuid "github.com/satori/go.uuid"
)


type Auth struct {
	ID				uuid.UUID				`gorm:"primaryKey" json:"id"`
	Username		string					`json:"username"`
	Password		string					`json:"-"`
	Role			int						`json:"role"`
	Name			string					`json:"name"`
	Nim				string					`json:"nim"`
	Email			string					`json:"email"`
	Address			string					`json:"address"`
	IsVerify		int						`json:"is_verify"`
	VerifyCode		int						`json:"verify_code"`		
	Rents			[]rentsModel.Rent		`gorm:"foreignKey:AuthID" json:"rents"`
	Reviews			[]reviewModel.Review	`gorm:"foreignKey:AuthReviewRefer" json:"reviews"`
	CreatedAt		time.Time				`json:"created_at"`
	UpdatedAt 		time.Time				`json:"updated_at"`
}