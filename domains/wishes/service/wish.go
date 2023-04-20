package service

import (
	"context"
	"fmt"

	"github.com/outcatcher/anwil/domains/wishes/service/schema"
)

// CreateWish creates new wish in the list.
func (s *service) CreateWish(
	ctx context.Context, wishlistUUID, description string, position int32,
) (string, error) {
	newID, err := s.storage.InsertWish(ctx, wishlistUUID, description, position)
	if err != nil {
		return "", fmt.Errorf("error creating wish: %w", err)
	}

	return newID, nil
}

// ListWishes return list of wishes in the wishlist.
func (s *service) ListWishes(ctx context.Context, wishlistUUID string) ([]schema.Wish, error) {
	wishes, err := s.storage.ListWishes(ctx, wishlistUUID)
	if err != nil {
		return nil, fmt.Errorf("error listing wishes: %w", err)
	}

	result := make([]schema.Wish, len(wishes))

	for i, wish := range wishes {
		result[i] = schema.Wish{
			Description: wish.Description,
			Fulfilled:   wish.Fulfilled,
		}
	}

	return result, nil
}
