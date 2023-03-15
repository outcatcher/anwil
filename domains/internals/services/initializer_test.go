package services

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	svcDTO "github.com/outcatcher/anwil/domains/internals/services/schema"
	"github.com/stretchr/testify/require"
)

type testService struct {
	id svcDTO.ServiceID

	expectedInitError error
	initCalled        atomic.Int32

	dependencies []svcDTO.ServiceID
}

// ID returns UUID unique to each testService instance.
func (t *testService) ID() svcDTO.ServiceID {
	if t.id == "" {
		t.id = svcDTO.ServiceID(uuid.New().String())
	}

	return t.id
}

// Init returns expectedInitError.
func (t *testService) Init(context.Context, interface{}) error {
	t.initCalled.Add(1)

	return t.expectedInitError
}

// DependsOn returns dependencies.
func (t *testService) DependsOn() []svcDTO.ServiceID {
	return t.dependencies
}

func svcMapping(services ...svcDTO.Service) svcDTO.ServiceMapping {
	mapping := make(svcDTO.ServiceMapping)

	for _, svc := range services {
		mapping[svc.ID()] = svc
	}

	return mapping
}

func TestInitialize(t *testing.T) { //nolint:funlen // this is a grouping test function
	t.Parallel()

	var emptyState any // not using state anyway

	t.Run("normal", func(t *testing.T) {
		t.Parallel()
		// normal
		svc1 := &testService{}
		svc2 := &testService{dependencies: []svcDTO.ServiceID{svc1.ID()}}

		mapping, err := Initialize(context.Background(), emptyState, svcMapping(svc2, svc1))
		require.NoError(t, err)
		require.NotNil(t, mapping)
	})

	t.Run("already initialized", func(t *testing.T) {
		t.Parallel()

		svc1 := &testService{}

		services := make([]svcDTO.Service, 1, 11)
		services[0] = svc1

		for i := 0; i < 10; i++ {
			services = append(services, &testService{dependencies: []svcDTO.ServiceID{svc1.ID()}})
		}

		_, err := Initialize(context.Background(), emptyState, svcMapping(services...))
		require.NoError(t, err)

		require.EqualValues(t, 1, svc1.initCalled.Load())
	})

	t.Run("cyclic direct", func(t *testing.T) {
		t.Parallel()

		id2 := svcDTO.ServiceID(uuid.New().String())

		svc1 := &testService{dependencies: []svcDTO.ServiceID{id2}}
		svc2 := &testService{id: id2, dependencies: []svcDTO.ServiceID{svc1.ID()}}

		_, err := Initialize(context.Background(), emptyState, svcMapping(svc2, svc1))
		require.ErrorIs(t, err, errCyclicServiceDependency)
	})

	t.Run("cyclic indirect", func(t *testing.T) {
		t.Parallel()

		id3 := svcDTO.ServiceID(uuid.New().String())

		svc1 := &testService{dependencies: []svcDTO.ServiceID{id3}}
		svc2 := &testService{dependencies: []svcDTO.ServiceID{svc1.ID()}}
		svc3 := &testService{dependencies: []svcDTO.ServiceID{svc2.ID()}, id: id3}

		_, err := Initialize(context.Background(), emptyState, svcMapping(svc2, svc1, svc3))
		require.ErrorIs(t, err, errCyclicServiceDependency)
	})

	t.Run("cyclic self", func(t *testing.T) {
		t.Parallel()

		id1 := svcDTO.ServiceID(uuid.New().String())
		svc1 := &testService{id: id1, dependencies: []svcDTO.ServiceID{id1}}

		_, err := Initialize(context.Background(), emptyState, svcMapping(svc1))
		require.ErrorIs(t, err, errCyclicServiceDependency)
	})

	t.Run("init error", func(t *testing.T) {
		t.Parallel()

		errInitFail := errors.New("error initializing testService") //nolint:goerr113 // this is intentional

		svc1 := &testService{expectedInitError: errInitFail}

		_, err := Initialize(context.Background(), emptyState, svcMapping(svc1))
		require.ErrorIs(t, err, errInitFail)
	})
}

func TestInitializeServiceWith(t *testing.T) {
	t.Parallel()

	var emptyState interface{}

	t.Run("normal", func(t *testing.T) {
		t.Parallel()

		svc1 := new(testService)

		unexpectedID := svc1.ID()
		expectedID := uuid.New().String()

		// initialize with a simple mutator
		err := InjectServiceWith(svc1, emptyState, func(service, state interface{}) error {
			tServ, ok := service.(*testService)
			require.True(t, ok)

			tServ.id = svcDTO.ServiceID(expectedID)

			return nil
		})

		require.NoError(t, err)
		require.NotEqualValues(t, expectedID, unexpectedID)
		require.EqualValues(t, expectedID, svc1.ID())
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		svc1 := new(testService)
		expectedErr := errors.New("expected svc init error") //nolint:goerr113

		// initialize with a simple mutator
		err := InjectServiceWith(svc1, emptyState, func(service, state interface{}) error {
			return expectedErr
		})

		require.ErrorIs(t, err, expectedErr)
	})
}
