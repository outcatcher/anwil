package services

import (
	"testing"

	svcSchema "github.com/outcatcher/anwil/domains/core/services/schema"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetServiceFromProvider(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		expected := new(testService)

		state := new(testState)
		state.
			On("Service", mock.AnythingOfType("ServiceID")).
			Return(expected)

		result, err := GetServiceFromProvider[*testService](state, "")
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})

	t.Run("nil service", func(t *testing.T) {
		t.Parallel()

		state := new(testState)
		state.
			On("Service", mock.AnythingOfType("ServiceID")).
			Return(nil)

		result, err := GetServiceFromProvider[*testService](state, "")
		require.ErrorIs(t, err, svcSchema.ErrMissingService)
		require.Nil(t, result)
	})

	t.Run("type mismatch", func(t *testing.T) {
		t.Parallel()

		state := new(testState)
		state.
			On("Service", mock.AnythingOfType("ServiceID")).
			Return(2)

		result, err := GetServiceFromProvider[*testService](state, "")
		require.ErrorIs(t, err, svcSchema.ErrInvalidType)
		require.Nil(t, result)
	})
}
