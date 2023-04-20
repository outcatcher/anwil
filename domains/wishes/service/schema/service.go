/*
Package schema contains service definition for Wishes service
*/
package schema

import (
	"context"

	svcSchema "github.com/outcatcher/anwil/domains/core/services/schema"
)

// ServiceWishes - ID for wishes service.
const ServiceWishes svcSchema.ServiceID = "wishes"

// WishesService defines the service of wishes.
type WishesService interface {
	svcSchema.Service

	CreateWishlist(ctx context.Context, userUUID, name string, visibility Visibility) error
	ListWishlists(ctx context.Context, userUUID string) ([]Wishlist, error)

	CreateWish(ctx context.Context, wishlistUUID, description string, position int32) (string, error)
	ListWishes(ctx context.Context, wishlistUUID string) ([]Wish, error)
}
