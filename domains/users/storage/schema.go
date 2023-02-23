package storage

// Wisher - entity of `wishers` table.
type Wisher struct {
	UUID     string `db:"uuid"`
	Username string `db:"username"`
	Password string `db:"password"`
	FullName string `db:"full_name"`
	Role     string `db:"role"`
	Enabled  bool   `db:"enabled"`
}
