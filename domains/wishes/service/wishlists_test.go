package service

import (
	"context"
	"database/sql"
	"testing"

	svcSchema "github.com/outcatcher/anwil/domains/core/services/schema"
	th "github.com/outcatcher/anwil/domains/core/testhelpers"
	"github.com/outcatcher/anwil/domains/wishes/service/schema"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	wishStorage "github.com/outcatcher/anwil/domains/wishes/storage"
)

type WishesSuite struct {
	suite.Suite
}

func (s *WishesSuite) TestService_CreateWishlist() {
	t := s.T()

	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		mockDB := &th.MockDBExecutor{}
		wishService := service{
			storage: wishStorage.New(mockDB),
		}

		ctx := context.Background()
		mockDB.
			On("GetContext", ctx, mock.Anything, mock.Anything, mock.Anything).
			Return(sql.ErrNoRows)
		mockDB.
			On("NamedExecContext", ctx, mock.Anything, mock.Anything).
			Return(th.MockSQLResult, nil)

		err := wishService.CreateWishlist(
			ctx,
			th.RandomString("uuid-", 10),
			th.RandomString("name-", 5),
			schema.VisibilityPublic,
		)
		require.NoError(t, err)
	})

	t.Run("existing", func(t *testing.T) {
		t.Parallel()

		mockDB := &th.MockDBExecutor{}
		wishService := service{
			storage: wishStorage.New(mockDB),
		}

		ctx := context.Background()
		mockDB.
			On("GetContext", ctx, mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		err := wishService.CreateWishlist(
			ctx,
			th.RandomString("uuid-", 10),
			th.RandomString("name-", 5),
			schema.VisibilityPublic,
		)
		require.ErrorIs(t, err, svcSchema.ErrConflict)
	})
}

func (s *WishesSuite) TestService_DeleteWishlist() {
	t := s.T()

	t.Parallel()

	t.Run("existing", func(t *testing.T) {
		t.Parallel()

		mockDB := &th.MockDBExecutor{}
		wishService := service{
			storage: wishStorage.New(mockDB),
		}

		ctx := context.Background()
		mockDB.
			On("ExecContext", ctx, mock.Anything, mock.Anything, mock.Anything).
			Return(th.MockSQLResult, nil)

		err := wishService.DeleteWishlist(
			ctx,
			th.RandomString("uuid-", 10),
			th.RandomString("name-", 5),
		)
		require.NoError(t, err)
	})
}

func TestService(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(WishesSuite))
}
