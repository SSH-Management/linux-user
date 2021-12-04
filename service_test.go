package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	mocks "github.com/SSH-Management/linux-user/__mocks__"
	"github.com/SSH-Management/linux-user/tests"
)

func TestCreateUser_Success(t *testing.T) {
	t.Parallel()

	dto := &mocks.User{}
	t.Cleanup(func() {
		tests.DeleteUser(t, dto.GetUsername())
	})

	assert := require.New(t)

	u := &UnixService{}

	linuxUser, err := u.Create(context.Background(), dto)

	assert.NoError(err)
	assert.Equal("/home/"+dto.GetUsername(), linuxUser.HomeFolder)
	assert.Equal(dto.GetUsername(), linuxUser.Username)
	assert.GreaterOrEqual(linuxUser.UserId, 1000)
}

func TestCreateUser_AlreadyExists(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	strErr := fmt.Sprintf("error while running the create user command, Exit Status: %d, Error: %v", 9, ErrUserAlreadyExists)

	logger := mocks.NewLogger()
	dto := &mocks.User{}

	t.Cleanup(func() {
		tests.DeleteUser(t, dto.GetUsername())
	})

	u := &UnixService{
		logger: logger,
	}

	tests.CreateUser(t, dto.GetUsername())

	_, err := u.Create(context.Background(), dto)

	assert.Error(err)
	assert.ErrorIs(err, ErrUserAlreadyExists)

	assert.Equal(logger.Len(), 1)
	assert.Equal(strErr, logger.Data()[0])
}

func TestDeleteUser_Success(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	logger := mocks.NewLogger()
	dto := &mocks.User{}
	u := &UnixService{
		logger: logger,
	}

	tests.CreateUser(t, dto.GetUsername())

	assert.NoError(u.Delete(context.Background(), dto.GetUsername()))
}

func TestDeleteUser_UserDoesNotExist(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	strErr := fmt.Sprintf("error while running the delete user command, Exit Status: %d, Error: %v", 6, ErrUserDoesNotExist)

	logger := mocks.NewLogger()
	dto := &mocks.User{}
	u := &UnixService{
		logger: logger,
	}

	err := u.Delete(context.Background(), dto.GetUsername())

	assert.Error(err)
	assert.ErrorIs(err, ErrUserDoesNotExist)
	assert.Equal(logger.Len(), 1)
	assert.Equal(strErr, logger.Data()[0])
}

func TestFindUserByUsername(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	logger := mocks.NewLogger()
	dto := &mocks.User{}

	t.Cleanup(func() {
		tests.DeleteUser(t, dto.GetUsername())
	})

	u := &UnixService{
		logger: logger,
	}

	tests.CreateUser(t, dto.GetUsername())

	linuxUser, err := u.Find(context.Background(), dto.GetUsername())

	assert.NoError(err)
	assert.NotEmpty(linuxUser)
	assert.Equal(dto.GetUsername(), linuxUser.Username)
	assert.GreaterOrEqual(linuxUser.UserId, 1000)
	assert.GreaterOrEqual(linuxUser.GroupId, 1000)
}



func TestFindUserByUsername_NotFound(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	logger := mocks.NewLogger()
	dto := &mocks.User{}

	u := &UnixService{
		logger: logger,
	}

	linuxUser, err := u.Find(context.Background(), dto.GetUsername())

	assert.Error(err)
	assert.Empty(linuxUser)
	assert.ErrorIs(err, ErrUserDoesNotExist)
}
