package entities

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	jwt.StandardClaims
	Id       int64  `json:"id"`
	RoleId   int32  `json:"role_id"`
	Username string `json:"username"`
}
