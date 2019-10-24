package models

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct{
	Name  string `json:"name"`
	Roles  string `json:"roles"`
	jwt.StandardClaims	
}