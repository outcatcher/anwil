package services

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	svcDTO "github.com/outcatcher/anwil/domains/services/dto"
	"github.com/stretchr/testify/require"
)

type testService struct {
	id                svcDTO.ServiceID
	expectedInitError error

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

func TestInitialize(t *testing.T) {
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
		svc2 := &testService{dependencies: []svcDTO.ServiceID{svc1.ID()}}
		svc3 := &testService{dependencies: []svcDTO.ServiceID{svc1.ID()}}

		mapping, err := Initialize(context.Background(), emptyState, svcMapping(svc1, svc2, svc3))
		require.NoError(t, err)
		require.NotNil(t, mapping)
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

		errInitFail := errors.New("error initializing testService")

		svc1 := &testService{expectedInitError: errInitFail}

		_, err := Initialize(context.Background(), emptyState, svcMapping(svc1))
		require.ErrorIs(t, err, errInitFail)
	})
}
