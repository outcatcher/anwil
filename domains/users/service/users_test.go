package service

import (
	"context"
	"crypto/ed25519"
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"testing"

	services "github.com/outcatcher/anwil/domains/core/services/schema"
	th "github.com/outcatcher/anwil/domains/core/testhelpers"
	"github.com/outcatcher/anwil/domains/users/service/schema"
	userStorage "github.com/outcatcher/anwil/domains/users/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	privateKey = "e3de69265ea200c17633b8b7ba90c17c15e96f3f1d0ad608d9f628e515c7e53b" +
		"d6507afe638ea0565709842d869581edfc5e5b6186a8215f6bed2504991ff9fb"
)

type UsersSuite struct {
	suite.Suite

	privateKey ed25519.PrivateKey
}

func (s *UsersSuite) SetupSuite() {
	t := s.T()

	key, err := hex.DecodeString(privateKey)
	require.NoError(t, err)

	s.privateKey = key
}

func TestUsers(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(UsersSuite))
}

func (s *UsersSuite) requireEqualPasswords(raw, encrypted string) {
	s.T().Helper()

	require.NoError(s.T(), validatePassword(raw, encrypted, s.privateKey))
}

func (s *UsersSuite) newService(mockDB *th.MockDBExecutor) *service {
	return &service{
		storage:    userStorage.New(mockDB),
		privateKey: s.privateKey,
	}
}

func setWisher(argIndex int, expectedUser userStorage.Wisher) func(args mock.Arguments) {
	return func(args mock.Arguments) {
		*(args.Get(argIndex).(*userStorage.Wisher)) = expectedUser //nolint:forcetypeassert
	}
}

func (s *UsersSuite) TestUsers_GetUser() {
	t := s.T()
	ctx := context.Background()

	t.Parallel()

	t.Run("no user", func(t *testing.T) {
		t.Parallel()

		mockDB := &th.MockDBExecutor{}
		mockDB.
			On("GetContext",
				ctx, new(userStorage.Wisher), mock.AnythingOfType("string"), mock.Anything,
			).
			Return(sql.ErrNoRows)

		users := s.newService(mockDB)

		_, err := users.GetUser(ctx, th.RandomString("user", 10))
		require.ErrorIs(t, err, services.ErrNotFound)
	})

	t.Run("existing", func(t *testing.T) {
		t.Parallel()

		expectedUser := userStorage.Wisher{
			Username: th.RandomString("user", 10),
			FullName: th.RandomString("full", 10),
		}

		mockDB := &th.MockDBExecutor{}
		mockDB.
			On("GetContext",
				ctx, new(userStorage.Wisher), mock.AnythingOfType("string"), []interface{}{expectedUser.Username},
			).
			Run(setWisher(1, expectedUser)).
			Return(nil)

		users := s.newService(mockDB)

		user, err := users.GetUser(ctx, expectedUser.Username)
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

		result := driver.RowsAffected(1)

		createdUser := &userStorage.Wisher{}

		mockDB := &th.MockDBExecutor{}
		mockDB.
			On("GetContext",
				ctx, new(userStorage.Wisher), mock.AnythingOfType("string"), mock.Anything,
			).
			Return(sql.ErrNoRows)
		mockDB.
			On("NamedExecContext",
				ctx, mock.AnythingOfType("string"), mock.AnythingOfType("storage.Wisher"),
			).
			Run(func(args mock.Arguments) {
				*createdUser = args.Get(2).(userStorage.Wisher) //nolint:forcetypeassert
			}).
			Return(result, nil)

		err := s.newService(mockDB).SaveUser(ctx, expectedUser)
		require.NoError(t, err)

		s.requireEqualPasswords(expectedUser.Password, createdUser.Password)
	})

	t.Run("existing user", func(t *testing.T) {
		t.Parallel()

		mockDB := &th.MockDBExecutor{}
		mockDB.
			On("GetContext",
				ctx, new(userStorage.Wisher), mock.AnythingOfType("string"), mock.Anything,
			).
			Return(nil)

		err := s.newService(mockDB).SaveUser(ctx, schema.User{
			Username: th.RandomString("username", 20),
			Password: th.RandomString("pwd", 20),
		})
		require.ErrorIs(t, err, services.ErrConflict)
	})
}

func (s *UsersSuite) TestUsers_GenerateUserToken() { //nolint:funlen
	t := s.T()
	ctx := context.Background()

	t.Parallel()

	rawPassword := th.RandomString("pwd-", 20)
	password, err := encrypt(rawPassword, s.privateKey)
	require.NoError(t, err)

	expectedUser := userStorage.Wisher{
		Username: th.RandomString("usr-", 5),
		Password: password,
		FullName: th.RandomString("Name ", 10),
	}

	t.Run("invalid user", func(t *testing.T) {
		t.Parallel()

		mockDB := &th.MockDBExecutor{}
		mockDB.
			On("GetContext",
				ctx, new(userStorage.Wisher), mock.AnythingOfType("string"), mock.Anything,
			).
			Return(sql.ErrNoRows)

		users := s.newService(mockDB)

		token, err := users.GenerateUserToken(ctx, schema.User{})
		require.ErrorIs(t, err, services.ErrNotFound)
		require.Empty(t, token)
	})

	t.Run("invalid password", func(t *testing.T) {
		t.Parallel()

		mockDB := &th.MockDBExecutor{}
		mockDB.
			On("GetContext",
				ctx, new(userStorage.Wisher), mock.AnythingOfType("string"), mock.Anything,
			).
			Run(setWisher(1, expectedUser)).
			Return(nil)

		testUser := schema.User{
			Username: expectedUser.Username,
			Password: "qwertyui",
			FullName: expectedUser.FullName,
		}

		token, err := s.newService(mockDB).GenerateUserToken(ctx, testUser)
		require.ErrorIs(t, err, services.ErrUnauthorized)
		require.Empty(t, token)
	})

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		mockDB := &th.MockDBExecutor{}
		mockDB.
			On("GetContext",
				ctx, new(userStorage.Wisher), mock.AnythingOfType("string"), mock.Anything,
			).
			Run(setWisher(1, expectedUser)).
			Return(nil)

		testUser := schema.User{
			Username: expectedUser.Username,
			Password: rawPassword,
			FullName: expectedUser.FullName,
		}

		token, err := s.newService(mockDB).GenerateUserToken(ctx, testUser)
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})
}
