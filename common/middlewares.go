package common

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	uuid "github.com/satori/go.uuid"
)

type JwtCustomClaims struct {
	ID				uuid.UUID	`gorm:"primaryKey" json:"id"`
	Username		string		`json:"username"`
	Role			int			`json:"role"`
	jwt.StandardClaims
}

func JwtMiddleware() echo.MiddlewareFunc {
	key:="SECRETKEY"
	return middleware.JWTWithConfig(middleware.JWTConfig{
		Claims: &JwtCustomClaims{},
		SigningKey: []byte(key),
	})
}

