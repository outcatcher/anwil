/*
Package handlers defines and populates API endpoints.
*/
package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/domains/api/middlewares"
	authDTO "github.com/outcatcher/anwil/domains/auth/dto"
	authHandlers "github.com/outcatcher/anwil/domains/auth/handlers"
	configDTO "github.com/outcatcher/anwil/domains/config/dto"
	logDTO "github.com/outcatcher/anwil/domains/logging/dto"
	services "github.com/outcatcher/anwil/domains/services/dto"
	usersDTO "github.com/outcatcher/anwil/domains/users/dto"
	userHandlers "github.com/outcatcher/anwil/domains/users/handlers"
)

func handleStatic(engine *gin.Engine, basePath string) {
	engine.Static("/static", basePath)
	engine.LoadHTMLGlob(filepath.Join(basePath, "*"))
}

type handlersState interface {
	logDTO.WithLogger
	usersDTO.WithUsers
	authDTO.WithAuth
	configDTO.WithConfig
}

type handlers struct {
	state handlersState

	baseGroup *gin.RouterGroup
	secGroup  *gin.RouterGroup
}

func newHandlers(state handlersState, engine *gin.Engine, baseAPIPath string) *handlers {
	h := &handlers{state: state}

	h.baseGroup = engine.Group(baseAPIPath, middlewares.ConvertErrors, middlewares.RequireJSON)

	h.secGroup = h.baseGroup.Group("/", middlewares.JWTAuth(state))

	return h
}

func (h *handlers) populate(funcs map[string]services.AddHandlersFunc) error {
	for name, hFunc := range funcs {
		if err := hFunc(h.baseGroup, h.secGroup); err != nil {
			return fmt.Errorf("error adding handlers for service %s: %w", name, err)
		}
	}

	return nil
}

func (h *handlers) populateCommon() {
	h.baseGroup.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/api/v1/swagger")
	})
	h.baseGroup.GET("/swagger", handleAPISpec)
	h.baseGroup.GET("/echo", handleEcho)

	h.secGroup.GET("/auth-echo", handleEcho)
}

// PopulateEndpoints populates endpoints for API.
func PopulateEndpoints(engine *gin.Engine, state handlersState) error {
	handleStatic(engine, state.Config().API.StaticPath)

	apiHandlers := newHandlers(state, engine, "/api/v1")

	apiHandlers.populateCommon()

	err := apiHandlers.populate(
		map[string]services.AddHandlersFunc{
			"auth":  authHandlers.AddAuthHandlers(state),
			"users": userHandlers.AddUserHandlers(state),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
