package api

import (
	"errors"
	"fmt"
)

var errAlreadyInitialized = errors.New("API already initialized")

func (s *State) initServices() error {
	if s.serviceMapping != nil {
		return fmt.Errorf("error initializing services: %w", errAlreadyInitialized)
	}

	s.serviceMapping = make(map[serviceID]any)

	for id, svc := range preparedServices {
		if err := svc.Init(s); err != nil {
			return fmt.Errorf("error initializing service %s: %w", id, err)
		}

		s.serviceMappingLock.Lock()
		s.serviceMapping[id] = svc
		s.serviceMappingLock.Unlock()
	}

	return nil
}
