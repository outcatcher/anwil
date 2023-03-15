/*
Package dto contains DTOs for auth service
*/
package dto

import "github.com/golang-jwt/jwt/v4"

// Claims - JWT payload contents.
type Claims struct {
	jwt.RegisteredClaims

	Username string `json:"username"`
}
