package schema

type Visibility string

const (
	VisibilityPrivate    Visibility = "private"
	VisibilityDirectLink Visibility = "direct_link"
	VisibilityPublic     Visibility = "public"
)

// Wish contains wish data.
type Wish struct {
	Description string // Wish text
	Fulfilled   bool   // Is wish marked as done
}

// Wishlist contains all wishlist data.
type Wishlist struct {
	Name       string
	Visibility Visibility
	Wishes     []Wish
}
