/*
Package storage contains db-related operations with users.
*/
package storage

import (
	"context"
	"fmt"

	storageSchema "github.com/outcatcher/anwil/domains/storage/schema"
)

// Storage - storage of wishes.
type Storage struct {
	db storageSchema.QueryExecutor
}

// New create a new UserStorage instance.
func New(db storageSchema.QueryExecutor) *Storage {
	return &Storage{db: db}
}

// GetWishlist returns wishlist by name and user UUID.
func (w *Storage) GetWishlist(ctx context.Context, userUUID, name string) (*Wishlist, error) {
	result := new(Wishlist)

	err := w.db.GetContext(
		ctx,
		result,
		`SELECT * FROM wishlists WHERE wisher_uuid = $1 AND name = $2 LIMIT 1;`,
		userUUID, name,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting wishlist: %w", err)
	}

	return result, nil
}

// SelectWishlists returns all wishlists for user.
func (w *Storage) SelectWishlists(ctx context.Context, userUUID string) ([]Wishlist, error) {
	result := make([]Wishlist, 0)

	err := w.db.SelectContext(
		ctx,
		&result,
		`SELECT * FROM wishlists WHERE wisher_uuid = $1 ORDER BY position;`,
		userUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("error selecting wishlists: %w", err)
	}

	return result, nil
}

// InsertWishlist creates a wish list.
func (w *Storage) InsertWishlist(ctx context.Context, data Wishlist) error {
	_, err := w.db.NamedExecContext(
		ctx,
		`INSERT INTO wishlists (wisher_uuid, name, visibility, position)
		VALUES (:wisher_uuid,
				:name,
				:visibility,
				(SELECT case when max(position) is NULL then 0 else max(position)+1 end
		 		FROM wishlists
		 		WHERE wisher_uuid = :wisher_uuid));`,
		data,
	)
	if err != nil {
		return fmt.Errorf("inserting wishlist failed: %w", err)
	}

	return nil
}

// DeleteWishlist removes existing wishlist.
func (w *Storage) DeleteWishlist(ctx context.Context, wisherID, name string) error {
	_, err := w.db.ExecContext(
		ctx,
		`DELETE FROM wishlists WHERE wisher_uuid = $1 AND name = $2;`,
		wisherID, name,
	)
	if err != nil {
		return fmt.Errorf("deleting wishlist failed: %w", err)
	}

	return nil
}

// InsertWish creates new wish in the wishlist.
func (w *Storage) InsertWish(
	ctx context.Context, wishlistUUID, description string, position int32,
) (string, error) {
	var uuid string

	err := w.db.GetContext(
		ctx,
		&uuid,
		`INSERT INTO wishes(wishlist_uuid, description, position) VALUES ($1, $2, $3) RETURNING uuid;`,
		wishlistUUID, description, position,
	)
	if err != nil {
		return "", fmt.Errorf("error inserting new wish: %w", err)
	}

	return uuid, nil
}

// ListWishes returns list of wishes in the wishlist.
func (w *Storage) ListWishes(ctx context.Context, wishlistUUID string) ([]Wish, error) {
	result := make([]Wish, 0)

	err := w.db.SelectContext(
		ctx,
		&result,
		`SELECT * FROM wishes WHERE wishlist_uuid = $1 ORDER BY position;`,
		wishlistUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("error selecting wishes: %w", err)
	}

	return result, nil
}

// UpdateWish updates wish attributes.
func (w *Storage) UpdateWish(ctx context.Context, wish Wish) error {
	_, err := w.db.NamedExecContext(
		ctx,
		`UPDATE wishes SET description = :description, fulfilled = :fulfilled WHERE uuid = :uuid;`,
		wish,
	)
	if err != nil {
		return fmt.Errorf("error updating wish: %w", err)
	}

	return nil
}
