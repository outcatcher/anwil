package api

import (
	"context"
	"errors"
	"fmt"
)

var errAlreadyInitialized = errors.New("API already initialized")

func (s *State) initServices(ctx context.Context) error {
	if s.serviceMapping != nil {
		return fmt.Errorf("error initializing services: %w", errAlreadyInitialized)
	}

	s.serviceMapping = make(map[serviceID]any)

	for id, svc := range preparedServices {
		if err := svc.Init(ctx, s); err != nil {
			return fmt.Errorf("error initializing service %s: %w", id, err)
		}

		s.serviceMappingLock.Lock()
		s.serviceMapping[id] = svc
		s.serviceMappingLock.Unlock()
	}

	return nil
}
