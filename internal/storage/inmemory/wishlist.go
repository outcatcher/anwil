package inmemory

import (
	"context"
	"fmt"

	"github.com/outcatcher/anwil/internal/storage/inmemory/mapstorage"
	"github.com/outcatcher/anwil/internal/storage/schema"
)

func (p *db) GetWishlistByID(_ context.Context, id int64) (*schema.Wishlist, error) {
	value, ok := p.wishlists.Get(id)
	if !ok {
		return nil, fmt.Errorf("missing wishlist %d: %w", id, schema.ErrNotFound)
	}

	return &value, nil
}

func (p *db) SaveWishlist(_ context.Context, wishlist schema.Wishlist) (int64, error) {
	if wishlist.ID == 0 {
		wishlist.ID = mapstorage.GenerateRandomInt64ID(p.wishlists)
	}

	p.wishlists.Set(wishlist.ID, wishlist)

	return wishlist.ID, nil
}
