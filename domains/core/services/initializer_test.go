package services

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	svcSchema "github.com/outcatcher/anwil/domains/core/services/schema"
	th "github.com/outcatcher/anwil/domains/core/testhelpers"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testService struct {
	mock.Mock

	injected bool
}

// init mocks service init returning given error.
func (m *testService) init(ctx context.Context, state any) error {
	mockArgs := m.Called(ctx, state)

	return mockArgs.Error(0)
}

type testState struct {
	mock.Mock
}

// Service - mock method to match ProvidingServices interface.
func (ts *testState) Service(id svcSchema.ServiceID) any {
	args := ts.Called(id)

	return args.Get(0)
}

func testServiceInit(svc *testService) svcSchema.ServiceInitFunc {
	return func(ctx context.Context, state any) (any, error) {
		if svc == nil {
			svc = &testService{}
			svc.On("init", ctx, state).Return(nil)
		}

		err := svc.init(ctx, state)
		if err != nil {
			return nil, err
		}

		return svc, nil
	}
}

func testServiceDefinition(svc *testService, dependencies ...svcSchema.ServiceID) svcSchema.ServiceDefinition {
	return svcSchema.ServiceDefinition{
		ID:        svcSchema.ServiceID(th.RandomString("id-", 5)),
		Init:      testServiceInit(svc),
		DependsOn: dependencies,
	}
}

func TestInitialize(t *testing.T) {
	t.Parallel()

	var emptyState any // not using state anyway

	t.Run("normal", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		mocked := &testService{}
		mocked.On("init", ctx, emptyState).Return(nil)

		// normal
		svc1 := testServiceDefinition(mocked)
		svc2 := testServiceDefinition(nil, svc1.ID)

		mapping, err := Initialize(ctx, emptyState, svc1, svc2)
		require.NoError(t, err)
		require.NotNil(t, mapping)

		mocked.AssertNumberOfCalls(t, "init", 1)
	})

	t.Run("already initialized", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		mocked := &testService{}
		mocked.On("init", ctx, emptyState).Return(nil)

		svc1 := testServiceDefinition(mocked)

		services := make([]svcSchema.ServiceDefinition, 1, 11)
		services[0] = svc1

		for i := 0; i < 10; i++ {
			services = append(services, testServiceDefinition(nil, svc1.ID))
		}

		_, err := Initialize(context.Background(), emptyState, services...)
		require.NoError(t, err)

		mocked.AssertNumberOfCalls(t, "init", 1)
	})

	t.Run("cyclic direct", func(t *testing.T) {
		t.Parallel()

		id2 := svcSchema.ServiceID(uuid.New().String())

		svc1 := testServiceDefinition(nil, id2)
		svc2 := svcSchema.ServiceDefinition{
			ID:        id2,
			Init:      testServiceInit(nil),
			DependsOn: []svcSchema.ServiceID{svc1.ID},
		}

		_, err := Initialize(context.Background(), emptyState, svc2, svc1)
		require.ErrorIs(t, err, errCyclicServiceDependency)
	})

	t.Run("cyclic indirect", func(t *testing.T) {
		t.Parallel()

		id3 := svcSchema.ServiceID(uuid.New().String())

		svc1 := testServiceDefinition(nil, id3)
		svc2 := testServiceDefinition(nil, svc1.ID)
		svc3 := svcSchema.ServiceDefinition{
			ID:        id3,
			Init:      testServiceInit(nil),
			DependsOn: []svcSchema.ServiceID{svc2.ID},
		}

		_, err := Initialize(context.Background(), emptyState, svc2, svc3, svc1)
		require.ErrorIs(t, err, errCyclicServiceDependency)
	})

	t.Run("cyclic self", func(t *testing.T) {
		t.Parallel()

		id1 := svcSchema.ServiceID(uuid.New().String())
		svc1 := svcSchema.ServiceDefinition{
			ID:        id1,
			Init:      testServiceInit(nil),
			DependsOn: []svcSchema.ServiceID{id1},
		}

		_, err := Initialize(context.Background(), emptyState, svc1)
		require.ErrorIs(t, err, errCyclicServiceDependency)
	})

	t.Run("depends on missing", func(t *testing.T) {
		t.Parallel()

		id2 := svcSchema.ServiceID(uuid.New().String())

		svc1 := testServiceDefinition(nil, id2)

		_, err := Initialize(context.Background(), emptyState, svc1)
		require.ErrorIs(t, err, errDefinitionMissing)
	})

	t.Run("init error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		errInitFail := errors.New("error initializing testService") //nolint:goerr113 // this is intentional

		mocked := &testService{}
		mocked.
			On("init", ctx, emptyState).
			Return(errInitFail)

		svc1 := testServiceDefinition(mocked)

		_, err := Initialize(ctx, emptyState, svc1)
		require.ErrorIs(t, err, errInitFail)
	})
}

func TestInitializeServiceWith(t *testing.T) {
	t.Parallel()

	var emptyState any

	t.Run("normal", func(t *testing.T) {
		t.Parallel()

		svc1 := new(testService)

		// initialize with a simple mutator
		err := InjectServiceWith(svc1, emptyState, func(service, state any) error {
			tServ, ok := service.(*testService)
			require.True(t, ok)

			tServ.injected = true

			return nil
		})

		require.NoError(t, err)
		require.EqualValues(t, true, svc1.injected)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		svc1 := new(testService)
		expectedErr := errors.New("expected svc init error") //nolint:goerr113

		err := InjectServiceWith(svc1, emptyState, func(service, state any) error {
			return expectedErr
		})

		require.ErrorIs(t, err, expectedErr)
	})
}

func TestValidateArgInterfaces(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		serv, state, err := ValidateArgInterfaces[*testService, *testState](&testService{}, &testState{})

		require.NoError(t, err)
		require.IsType(t, new(testService), serv)
		require.IsType(t, new(testState), state)
	})

	t.Run("not providing", func(t *testing.T) {
		t.Parallel()

		_, _, err := ValidateArgInterfaces[*testService, *testState](&testService{}, &testService{})

		require.ErrorIs(t, err, errNotProvided)
	})

	t.Run("not required", func(t *testing.T) {
		t.Parallel()

		_, _, err := ValidateArgInterfaces[*testService, *testState](&testState{}, &testState{})

		require.ErrorIs(t, err, errNotNeeded)
	})
}
