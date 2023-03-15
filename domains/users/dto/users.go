/*
Package dto contains DTOs for User domain.
*/
package dto

// User holds user data.
type User struct {
	Username string `json:"username"`
	Password string `json:"-"` // hex-encoded password, make sure it's not reaching JSON
	FullName string `json:"full_name"`
}
