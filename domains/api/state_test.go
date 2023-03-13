package api

import (
	"bytes"
	"context"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/domains/logging"
	"github.com/stretchr/testify/require"
)

func TestLoggerFromCtx_ginCtx(t *testing.T) {
	t.Parallel()

	// real server workflow

	baseCtx := context.Background()
	expectWriter := new(bytes.Buffer)
	logger := log.New(expectWriter, "", 0)

	ctx := logging.CtxWithLogger(baseCtx, logger)
	require.NotNil(t, ctx)

	expectedString := "test me in context\n"

	engine := gin.New()

	engine.Handle(http.MethodGet, "/", func(c *gin.Context) {
		logger := logging.LoggerFromCtx(c.Request.Context())

		logger.Print(expectedString)

		c.AbortWithStatus(http.StatusOK)
	})

	server := &http.Server{ //nolint:exhaustruct
		Addr:              "localhost:7343",
		Handler:           engine,
		ReadHeaderTimeout: defaultTimeout,
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

	t.Log(expectWriter.String())
}
