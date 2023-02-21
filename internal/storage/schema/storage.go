package schema

import "context"

// User - data stored for user.
type User struct {
	Username string `json:"username"`
	Password string `json:"-"`
	FullName string `json:"full_name"`
}

// Wish - single wishlist entry.
type Wish struct {
	Private bool `json:"private"`
}

// Wishlist - wishlist data.
type Wishlist struct {
	ID     int64  `json:"id"`
	Owner  User   `json:"owner"`
	Wishes []Wish `json:"wishes"`
}

// Storage - хранилище данных.
type Storage interface {
	// GetUser returns existing user by username.
	GetUser(ctx context.Context, username string) (*User, error)
	// SaveUser creates new or updates existing user.
	SaveUser(ctx context.Context, user User) error
	GetWishlistByID(ctx context.Context, id int64) (*Wishlist, error)
	SaveWishlist(ctx context.Context, wishlist Wishlist) (int64, error)
}
