package requestmodels

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId      string `json:"user_id"`
	DormitoryId int    `json:"dormitory_id"`
	Type        string `json:"type"`

	jwt.RegisteredClaims
}
