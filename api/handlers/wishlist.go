package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/api/ctxhelpers"
	"github.com/outcatcher/anwil/internal/storage"
	"github.com/outcatcher/anwil/internal/storage/schema"
)

// sanitizeWishlist remove private data from the wishlist.
func sanitizeWishlist(src *schema.Wishlist) *schema.Wishlist {
	if src == nil {
		return nil
	}

	publicWishes := make([]schema.Wish, 0, len(src.Wishes))

	for _, wish := range src.Wishes {
		if !wish.Private {
			publicWishes = append(publicWishes, wish)
		}
	}

	return &schema.Wishlist{
		ID:     src.ID,
		Owner:  src.Owner, // nothing to hide ATM
		Wishes: publicWishes,
	}
}

func handleGetWishlist(c *gin.Context) {
	intID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		c.Abort()

		return
	}

	wishlist, err := storage.Storage().GetWishlistByID(c.Request.Context(), intID)
	if errors.Is(err, schema.ErrNotFound) {
		c.String(http.StatusNotFound, "no wishlist with ID %d found", intID)
		c.Abort()

		return
	}

	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	if wishlist == nil {
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	username := c.GetString(ctxhelpers.CtxKeyUsername)
	if wishlist.Owner.Username != username {
		wishlist = sanitizeWishlist(wishlist)
	}

	c.JSON(http.StatusOK, wishlist)
}
