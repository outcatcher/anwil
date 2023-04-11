package logging

import (
	"bytes"
	"context"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtxWithLogger(t *testing.T) {
	t.Parallel()

	baseCtx := context.Background()
	expectWriter := new(bytes.Buffer)
	logger := log.New(expectWriter, "", 0)

	ctx := CtxWithLogger(baseCtx, logger)
	require.NotNil(t, ctx)

	actualLogger := LoggerFromCtx(ctx)
	require.NotNil(t, actualLogger)
	require.Equal(t, logger, actualLogger)

	expectedText := "test me please\n"

	actualLogger.Print(expectedText)

	require.EqualValues(t, expectedText, expectWriter.String())
}

func TestLoggerFromCtx_Nil(t *testing.T) {
	t.Parallel()

	baseCtx := context.Background()

	actualLogger := LoggerFromCtx(baseCtx)
	require.Equal(t, io.Discard, actualLogger.Writer())
}
