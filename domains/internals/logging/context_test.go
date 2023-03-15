package logging

import (
	"bytes"
	"context"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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
	require.Nil(t, actualLogger)
}

func TestLoggerFromCtx_ginCtx(t *testing.T) {
	t.Parallel()

	// real server workflow

	baseCtx := context.Background()
	expectWriter := new(bytes.Buffer)
	logger := log.New(expectWriter, "", 0)

	ctx := CtxWithLogger(baseCtx, logger)
	require.NotNil(t, ctx)

	expectedString := "test me in context\n"

	engine := gin.New()

	engine.Handle(http.MethodGet, "/", func(c *gin.Context) {
		logger := LoggerFromCtx(c.Request.Context())

		logger.Print(expectedString)

		c.AbortWithStatus(http.StatusOK)
	})

	server := &http.Server{
		Addr:              "localhost:7343",
		Handler:           engine,
		ReadHeaderTimeout: time.Second,
		BaseContext:       func(_ net.Listener) context.Context { return ctx },
	}

	go func() {
		_ = server.ListenAndServe()
	}()

	t.Cleanup(func() {
		require.NoError(t, server.Shutdown(context.Background()))
	})

	time.Sleep(2 * time.Millisecond) // oof

	req, err := http.NewRequest(http.MethodGet, "http://"+server.Addr, nil)
	require.NoError(t, err)

	_, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.NotEmpty(t, expectWriter.String())
	t.Log(expectWriter.String())
}
