package services

import (
	"fmt"

	"github.com/outcatcher/anwil/domains/core/services/schema"
)

// GetServiceFromProvider returns service of exact type by ID.
func GetServiceFromProvider[T any](p schema.ProvidingServices, id schema.ServiceID) (T, error) {
	rawService := p.Service(id)

	var svc T

	if rawService == nil {
		return svc, fmt.Errorf("%w: %s", schema.ErrMissingService, id)
	}

	svc, ok := rawService.(T)
	if !ok {
		return svc, fmt.Errorf("%w: actual service type %T", schema.ErrInvalidType, rawService)
	}

	return svc, nil
}
