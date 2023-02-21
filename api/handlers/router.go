package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/api/middlewares"
	"github.com/outcatcher/anwil/config"
)

// NewRouter creates new GIN engine for Anwil API.
func NewRouter(basePath string, cfg *config.ServerConfiguration, middles ...gin.HandlerFunc) (*gin.Engine, error) {
	engine := gin.New()
	engine.Use(middles...)

	serveStaticHTML(engine, basePath)

	api, err := newAPI(cfg)
	if err != nil {
		return nil, err
	}

	api.populateAPIEndpoints(engine)

	return engine, nil
}

func serveStaticHTML(engine *gin.Engine, basePath string) {
	engine.Static("/static", filepath.Join(basePath, "static"))
	engine.LoadHTMLGlob(filepath.Join(basePath, "static", "*.html"))
}

func (s server) populateAPIEndpoints(engine *gin.Engine) {
	api := engine.Group("/api/v1")

	api.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/api/v1/swagger")
	})
	api.GET("/swagger", s.handleAPISpec)
	api.GET("/echo", s.handleEcho)

	api.POST("/token", handleAuthorize)

	secure := api.Group("/", middlewares.JWTAuth(s.PrivateKey))
	secure.GET("/auth-echo", s.handleEcho)

	secure.GET("/wishlist/:id", handleGetWishlist)
}
