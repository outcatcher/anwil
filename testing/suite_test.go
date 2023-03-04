//go:build acceptance

package testing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/domains/api"
	apiDTO "github.com/outcatcher/anwil/domains/api/dto"
	"github.com/outcatcher/anwil/domains/config"
	storageDTO "github.com/outcatcher/anwil/domains/config/dto"
	"github.com/outcatcher/anwil/domains/storage"
	usersDTO "github.com/outcatcher/anwil/domains/users/dto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	containerName           = "test_postgres"
	containerStartPeriod    = 500 * time.Millisecond
	containerStartUpTimeout = 10 * time.Second

	debugUsername = "debug"
	debugFullName = "Debug Wisher"
	debugPassword = "debug"
)

// AnwilSuite - handlers tests.
type AnwilSuite struct {
	suite.Suite

	apiHandler http.HandlerFunc
}

// requestJSON sends request with `content-type: application/json`.
func (s *AnwilSuite) requestJSON(
	method string, url *url.URL, body any, headers map[string]string,
) *httptest.ResponseRecorder {
	t := s.T()
	t.Helper()

	var reader io.Reader

	switch b := body.(type) {
	case nil:
	case string:
		reader = strings.NewReader(b)
	case []byte:
		reader = bytes.NewReader(b)
	case io.Reader:
		reader = b
	default:
		// convert to JSON
		data, err := json.Marshal(b)
		require.NoError(t, err)

		reader = bytes.NewReader(data)
	}

	if headers == nil {
		headers = make(map[string]string)
	}

	headers["content-type"] = gin.MIMEJSON

	return s.request(method, url, reader, headers)
}

func (s *AnwilSuite) request(
	method string, url *url.URL, body io.Reader, headers map[string]string,
) *httptest.ResponseRecorder {
	t := s.T()
	t.Helper()

	responseRecorder := httptest.NewRecorder()

	request := httptest.NewRequest(method, url.Path, body)
	request.URL.RawQuery = url.Query().Encode()

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	s.apiHandler(responseRecorder, request)

	return responseRecorder
}

func (s *AnwilSuite) SetupSuite() {
	t := s.T()
	ctx := context.Background()

	gin.SetMode(gin.ReleaseMode) // no need for request logs

	configPath := "./fixtures/test_config.yaml"

	cfg, err := config.LoadServerConfiguration(ctx, configPath)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	startDBContainer(t, ctx, cfg.DB)
	require.NoError(t, storage.ApplyMigrations(cfg.DB))

	apiState, err := api.Init(ctx, configPath)
	require.NoError(t, err)

	createDebugUser(t, ctx, apiState)

	router, err := apiState.NewRouter(ctx)
	require.NoError(t, err)

	// gin engine as a handler function
	s.apiHandler = func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r)
	}
}

func mapToSlice(src map[string]string) []string {
	result := make([]string, 0, len(src))

	for k, v := range src {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}

	return result
}

func startDBContainer(
	t *testing.T, ctx context.Context, dbConfig storageDTO.DatabaseConfiguration,
) {
	t.Helper()

	dockerClient, err := client.NewClientWithOpts()
	require.NoError(t, err)

	postgresPort := strconv.Itoa(dbConfig.Port)
	envMap := map[string]string{
		"POSTGRES_PORT":     postgresPort,
		"POSTGRES_PASSWORD": dbConfig.Password,
		"POSTGRES_USER":     dbConfig.Username,
		"POSTGRES_DB":       dbConfig.DatabaseName,
	}

	natPgPort, err := nat.NewPort("tcp", postgresPort)
	require.NoError(t, err)

	containerConfig := &container.Config{
		Hostname: containerName,
		ExposedPorts: nat.PortSet{
			natPgPort: struct{}{},
		},
		Env: mapToSlice(envMap),
		Cmd: []string{"-p", postgresPort},
		Healthcheck: &container.HealthConfig{
			Test: []string{
				"CMD-SHELL", fmt.Sprintf("nc %s %s -zv || exit 1", containerName, postgresPort),
			},
			Interval:    time.Second,
			StartPeriod: containerStartPeriod,
			Retries:     1,
		},
		Image: "postgres:15.2-alpine3.17",
	}

	created, err := dockerClient.ContainerCreate(
		ctx,
		containerConfig,
		&container.HostConfig{
			PortBindings: nat.PortMap{
				natPgPort: []nat.PortBinding{
					{
						HostPort: postgresPort,
					},
				},
			},
		},
		nil,
		nil,
		containerName,
	)
	require.NoError(t, err)

	err = dockerClient.ContainerStart(ctx, created.ID, types.ContainerStartOptions{})
	require.NoError(t, err)

	t.Cleanup(func() {
		cleanupCtx := context.Background()

		minute := time.Minute

		err = dockerClient.ContainerStop(cleanupCtx, created.ID, &minute)
		require.NoError(t, err)

		err = dockerClient.ContainerRemove(cleanupCtx, created.ID, types.ContainerRemoveOptions{
			Force: true,
		})
		require.NoError(t, err)
	})

	time.Sleep(containerStartPeriod) // wait for healthcheck to be working

	waitForDBUp(t, ctx, dockerClient, created.ID)
}

func waitForDBUp(t *testing.T, ctx context.Context, dockerClient *client.Client, containerID string) {
	t.Helper()

	require.NotNil(t, dockerClient)

	deadlineCtx, cancel := context.WithTimeout(ctx, containerStartUpTimeout)
	defer cancel()

	ready := false

	// wait for database to start
	for !ready {
		state, err := dockerClient.ContainerInspect(deadlineCtx, containerID)
		require.NoError(t, err)

		if state.State.Status != "running" {
			continue
		}

		if state.State.Health == nil {
			continue
		}

		if state.State.Health.Status == types.Healthy {
			ready = true
		}
	}
}

func createDebugUser(t *testing.T, ctx context.Context, state apiDTO.State) {
	err := state.Users().SaveUser(ctx, usersDTO.User{
		Username: debugUsername,
		Password: debugPassword,
		FullName: debugFullName,
	})
	require.NoError(t, err)
}

func TestRun(t *testing.T) {
	t.Parallel()

	testSuite := new(AnwilSuite)

	suite.Run(t, testSuite)
}

func parseRequestURL(t *testing.T, path string) *url.URL {
	t.Helper()

	requestURL, err := url.ParseRequestURI(path)
	require.NoError(t, err)

	return requestURL
}
