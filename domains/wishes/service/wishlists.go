package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	svcSchema "github.com/outcatcher/anwil/domains/core/services/schema"
	"github.com/outcatcher/anwil/domains/wishes/service/schema"
	"github.com/outcatcher/anwil/domains/wishes/storage"
)

var errCreateWishlist = errors.New("error creating wishlist")

// CreateWishlist creates new wish list for user.
func (s *service) CreateWishlist(ctx context.Context, userUUID, name string, visibility schema.Visibility) error {
	existing, err := s.storage.GetWishlist(ctx, userUUID, name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %w", errCreateWishlist, err)
	}

	if existing != nil {
		return fmt.Errorf("%w: %w", errCreateWishlist, svcSchema.ErrConflict)
	}

	err = s.storage.InsertWishlist(ctx, storage.Wishlist{
		WisherID:   userUUID,
		Name:       name,
		Visibility: string(visibility),
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errCreateWishlist, err)
	}

	return nil
}

// DeleteWishlist removes existing wishlist.
func (s *service) DeleteWishlist(ctx context.Context, userUUID, name string) error {
	err := s.storage.DeleteWishlist(ctx, userUUID, name)
	if err != nil {
		return fmt.Errorf("error deleting wishlist: %w", err)
	}

	return nil
}

// ListWishlists lists all user wishlists.
func (s *service) ListWishlists(ctx context.Context, userUUID string) ([]schema.Wishlist, error) {
	wishlists, err := s.storage.SelectWishlists(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("error listing wishlists: %w", err)
	}

	result := make([]schema.Wishlist, len(wishlists))

	for i, wishlist := range wishlists {
		wishes, err := s.ListWishes(ctx, wishlist.UUID)
		if err != nil {
			return nil, fmt.Errorf("error listing wishes for wishlist: %w", err)
		}

		result[i] = schema.Wishlist{
			Name:       wishlist.Name,
			Visibility: schema.Visibility(wishlist.Visibility),
			Wishes:     wishes,
		}
	}

	return result, nil
}
