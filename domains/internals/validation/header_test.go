package validation

import (
	"context"
	"fmt"
	"testing"

	v10 "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

const testToken = "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9." +
	"eyJ1c2VybmFtZSI6InJhbmRvbS11c2VybmFtZSJ9." +
	"nVxpOGoHA9cggY3yGY9RJdZLYdPYnkBTClfG5HTLtLA4uEUKz5tdjlKvGHr0DkT9AVq1tiaC1SxC1ICcV4wECg"

type testHeaders struct {
	TestToken string `header:"Test-Token" validate:"required,jwt-header"`
}

func TestValidateHeaderCtx_ok(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		value any
	}{
		{
			"normal",
			testHeaders{fmt.Sprintf("Bearer %s", testToken)},
		},
		{
			"pointer",
			&testHeaders{fmt.Sprintf("Bearer %s", testToken)},
		},
	}

	for _, data := range cases {
		data := data

		t.Run(data.name, func(t *testing.T) {
			t.Parallel()

			require.NoError(t, ValidateHeaderCtx(context.Background(), data.value))
		})
	}
}

func TestValidateHeaderCtx_validationErr(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		headers any
	}{
		{
			"no header",
			testHeaders{},
		},
		{
			"no bearer",
			testHeaders{testToken},
		},
		{
			"invalid token",
			testHeaders{"Bearer 1234"},
		},
		{
			"pointer",
			&testHeaders{"Bearer 1234"},
		},
	}

	for _, data := range cases {
		data := data
		ctx := context.Background()

		t.Run(data.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateHeaderCtx(ctx, data.headers)
			require.ErrorIs(t, err, ErrValidationFailed)
			require.ErrorContains(t, err, "Test-Token")
			t.Log(err)
		})
	}
}

func TestValidateHeaderCtx_notHeader(t *testing.T) {
	t.Parallel()

	value := testJSON{}
	ctx := context.Background()

	err := ValidateHeaderCtx(ctx, value)
	require.ErrorIs(t, err, ErrValidationFailed)
	t.Log(err)
}

func TestValidateHeaderCtx_invalidData(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	targetErr := new(v10.InvalidValidationError)

	err := ValidateHeaderCtx(ctx, nil)
	require.ErrorAs(t, err, &targetErr)
	t.Log(err)
}
