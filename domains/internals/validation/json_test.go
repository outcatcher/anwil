package validation

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type testJSON struct {
	Field string `json:"field" validate:"required"`
}

// There are two main differences with basic validator:
// 1. JSON tag value in error message
// 2. Different validation tag (`validate` instead of `binding`)

func TestValidateJSONCtx_validationErr(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		data      any
		fieldName string
	}{
		{"normal", testJSON{}, "field"},
		{"pointer", &testJSON{}, "field"},
		{"no-JSON", struct {
			Field string `validate:"required"`
		}{}, "Field"},
		{"inherited", struct {
			Field string `json:"field" validate:"min=2"`
		}{}, "field"},
	}

	for _, testData := range cases {
		testData := testData

		t.Run(testData.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			err := ValidateJSONCtx(ctx, testData.data)
			require.ErrorIs(t, err, ErrValidationFailed)
			require.ErrorContains(t, err, fmt.Sprintf("field '%s'", testData.fieldName))
			t.Log(err)
		})
	}
}

func TestValidateJSONCtx_ok(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		value any
	}{
		{
			"normal",
			testJSON{fmt.Sprintf("Bearer %s", testToken)},
		},
		{
			"pointer",
			&testJSON{fmt.Sprintf("Bearer %s", testToken)},
		},
	}

	for _, data := range cases {
		data := data

		t.Run(data.name, func(t *testing.T) {
			t.Parallel()

			require.NoError(t, ValidateJSONCtx(context.Background(), data.value))
		})
	}
}
