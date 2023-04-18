package service

import (
	"context"
	"encoding/hex"
	"testing"

	services "github.com/outcatcher/anwil/domains/core/services/schema"
	th "github.com/outcatcher/anwil/domains/core/testhelpers"
	"github.com/outcatcher/anwil/domains/users/service/schema"
	"github.com/outcatcher/anwil/domains/users/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	privateKey = "e3de69265ea200c17633b8b7ba90c17c15e96f3f1d0ad608d9f628e515c7e53b" +
		"d6507afe638ea0565709842d869581edfc5e5b6186a8215f6bed2504991ff9fb"
)

type UsersSuite struct {
	suite.Suite

	users service
}

func (s *UsersSuite) SetupSuite() {
	t := s.T()

	key, err := hex.DecodeString(privateKey)
	require.NoError(t, err)

	userService := service{
		storage:    storage.NewMock(),
		privateKey: key,
	}

	s.users = userService
}

func TestUsers(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(UsersSuite))
}

func (s *UsersSuite) requireEqualPasswords(raw, encrypted string) {
	s.T().Helper()

	require.NoError(s.T(), validatePassword(raw, encrypted, s.users.privateKey))
}

func (s *UsersSuite) TestUsers_GetUser() {
	t := s.T()
	ctx := context.Background()

	t.Parallel()

	t.Run("no user", func(t *testing.T) {
		t.Parallel()

		_, err := s.users.GetUser(ctx, th.RandomString("user", 10))
		require.ErrorIs(t, err, services.ErrNotFound)
	})

	t.Run("existing", func(t *testing.T) {
		t.Parallel()

		expectedUser := storage.Wisher{
			Username: th.RandomString("user", 10),
			FullName: th.RandomString("full", 10),
		}

		err := s.users.storage.InsertUser(ctx, expectedUser)
		require.NoError(t, err)

		user, err := s.users.GetUser(ctx, expectedUser.Username)
		require.NoError(t, err)

		require.EqualValues(t, expectedUser.Username, user.Username)
		require.EqualValues(t, expectedUser.FullName, user.FullName)
	})
}

func (s *UsersSuite) TestUsers_SaveUser() {
	t := s.T()
	ctx := context.Background()

	t.Parallel()

	t.Run("no user", func(t *testing.T) {
		t.Parallel()

		expectedUser := schema.User{
			Username: th.RandomString("user", 10),
			Password: th.RandomString("pwd", 10),
			FullName: th.RandomString("full", 10),
		}

		err := s.users.SaveUser(ctx, expectedUser)
		require.NoError(t, err)

		userInStorage, err := s.users.storage.GetUser(ctx, expectedUser.Username)
		require.NoError(t, err)

		require.EqualValues(t, expectedUser.Username, userInStorage.Username)
		require.EqualValues(t, expectedUser.FullName, userInStorage.FullName)
		s.requireEqualPasswords(expectedUser.Password, userInStorage.Password)
	})

	t.Run("existing user", func(t *testing.T) {
		t.Parallel()

		expectedUser := storage.Wisher{
			Username: th.RandomString("user", 10),
			Password: th.RandomString("pwd", 20),
		}

		err := s.users.storage.InsertUser(ctx, expectedUser)
		require.NoError(t, err)

		err = s.users.SaveUser(ctx, schema.User{
			Username: expectedUser.Username,
			Password: th.RandomString("pwd", 20),
		})
		require.ErrorIs(t, err, services.ErrConflict)
	})
}

func (s *UsersSuite) createTestUser(ctx context.Context) schema.User {
	t := s.T()
	t.Helper()

	testUser := schema.User{
		Username: th.RandomString("usr-", 5),
		Password: th.RandomString("pwd-", 20),
		FullName: "Test User",
	}

	err := s.users.SaveUser(ctx, testUser)
	require.NoError(t, err)

	return testUser
}

func (s *UsersSuite) TestUsers_GenerateUserToken() {
	t := s.T()
	ctx := context.Background()

	t.Parallel()

	t.Run("invalid user", func(t *testing.T) {
		t.Parallel()

		token, err := s.users.GenerateUserToken(ctx, schema.User{})
		require.ErrorIs(t, err, services.ErrNotFound)
		require.Empty(t, token)
	})

	t.Run("invalid password", func(t *testing.T) {
		t.Parallel()

		testUser := s.createTestUser(ctx)
		testUser.Password = "qwertyui"

		token, err := s.users.GenerateUserToken(ctx, testUser)
		require.ErrorIs(t, err, services.ErrUnauthorized)
		require.Empty(t, token)
	})

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		testUser := s.createTestUser(ctx)

		token, err := s.users.GenerateUserToken(ctx, testUser)
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})
}
