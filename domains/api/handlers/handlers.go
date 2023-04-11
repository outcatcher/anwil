/*
Package handlers defines and populates API endpoints.
*/
package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/outcatcher/anwil/domains/api/middlewares"
	configSchema "github.com/outcatcher/anwil/domains/core/config/schema"
	logSchema "github.com/outcatcher/anwil/domains/core/logging/schema"
	services "github.com/outcatcher/anwil/domains/core/services/schema"
	userHandlers "github.com/outcatcher/anwil/domains/users/handlers"
	"github.com/outcatcher/anwil/domains/users/service/schema"
)

func handleStatic(engine *fiber.App, basePath string) {
	engine.Static("/static", basePath)
}

type handlersState interface {
	logSchema.WithLogger
	schema.WithUsers
	configSchema.WithConfig
}

type handlers struct {
	state handlersState

	baseGroup fiber.Router
}

func newHandlers(state handlersState, engine *fiber.App, baseAPIPath string) *handlers {
	h := &handlers{state: state}

	h.baseGroup = engine.Group(baseAPIPath)

	return h
}

func (h *handlers) populate(funcs map[string]services.AddHandlersFunc) error {
	for name, hFunc := range funcs {
		if err := hFunc(h.baseGroup); err != nil {
			return fmt.Errorf("error adding handlers for service %s: %w", name, err)
		}
	}

	return nil
}

func (h *handlers) populateCommon(state handlersState) {
	h.baseGroup.Get("/echo", handleEcho)
	h.baseGroup.Get("/auth-echo", middlewares.JWTAuth(state), handleEcho)
}

// PopulateEndpoints populates endpoints for API.
func PopulateEndpoints(engine *fiber.App, state handlersState) error {
	handleStatic(engine, state.Config().API.StaticPath)

	apiHandlers := newHandlers(state, engine, "/api/v1")

	apiHandlers.populateCommon(state)

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
