/*
Package handlers defines and populates API endpoints.
*/
package handlers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/outcatcher/anwil/domains/api/middlewares"
	configSchema "github.com/outcatcher/anwil/domains/core/config/schema"
	logSchema "github.com/outcatcher/anwil/domains/core/logging/schema"
	services "github.com/outcatcher/anwil/domains/core/services/schema"
	userHandlers "github.com/outcatcher/anwil/domains/users/handlers"
	"github.com/outcatcher/anwil/domains/users/service/schema"
)

func handleStatic(engine *echo.Echo, basePath string) {
	engine.Static("/static", basePath)
}

type handlersState interface {
	logSchema.WithLogger
	schema.WithUsers
	configSchema.WithConfig
}

type handlers struct {
	state handlersState

	baseGroup *echo.Group
	secGroup  *echo.Group
}

func newHandlers(state handlersState, echo *echo.Echo, baseAPIPath string) *handlers {
	h := &handlers{state: state}

	h.baseGroup = echo.Group(baseAPIPath)
	h.secGroup = h.baseGroup.Group("", middlewares.JWTAuth(state))

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
	h.baseGroup.GET("/echo", handleEcho)
	h.secGroup.GET("/auth-echo", handleEcho)
}

// PopulateEndpoints populates endpoints for API.
func PopulateEndpoints(engine *echo.Echo, state handlersState) error {
	handleStatic(engine, state.Config().API.StaticPath)

	apiHandlers := newHandlers(state, engine, "/api/v1")

	apiHandlers.populateCommon()

	err := apiHandlers.populate(
		map[string]services.AddHandlersFunc{
			"users": userHandlers.AddUserHandlers(state),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
