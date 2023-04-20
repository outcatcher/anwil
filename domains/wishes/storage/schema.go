package storage

// Wishlist - entity in `wishlist` table.
type Wishlist struct {
	UUID       string `db:"uuid"`
	WisherID   string `db:"wisher_uuid"`
	Name       string `db:"name"`
	Visibility string `db:"visibility"`
	Position   int32  `db:"position"`
}

// Wish - entity of `wish` table.
type Wish struct {
	UUID         string `db:"uuid"`
	WishlistUUID string `db:"wishlist_uuid"`
	Description  string `db:"description"`
	Fulfilled    bool   `db:"fulfilled"`
}
