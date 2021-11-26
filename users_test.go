package user

import (
	"context"
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"

	mocks "github.com/SSH-Management/linux-user/__mocks__"
)

func deleteUser(t *testing.T, username string) {
	cmd := exec.CommandContext(
		context.Background(),
		DeleteCommand,
		"-f",
		"-r",
		username,
	)

	err := cmd.Run()
	if err != nil {
		t.Errorf("Error while deleting user: %v\n", err)
		t.FailNow()
	}
}

func createUser(t *testing.T, username string) {
	cmd := exec.CommandContext(
		context.Background(),
		AddCommand,
		"--create-home",
		"--password",
		"password",
		"--shell",
		"/bin/sh",
		"--user-group",
		username,
	)

	err := cmd.Run()
	if err != nil {
		t.Errorf("Error while creating user: %v\n", err)
		t.FailNow()
	}
}

func TestCreateUser_Success(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	dto := &mocks.User{}
	u := &UnixService{}

	defer deleteUser(t, dto.GetUsername())

	linuxUser, err := u.Create(context.Background(), dto)

	assert.NoError(err)
	assert.Equal("/home/"+dto.GetUsername(), linuxUser.HomeFolder)
	assert.GreaterOrEqual(linuxUser.UserId, 1000)
}

func TestCreateUser_AlreadyExists(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	strErr := fmt.Sprintf("error while running the create user command, Exit Status: %d, Error: %v", 9, ErrUserAlreadyExists)

	logger := mocks.NewLogger()
	dto := &mocks.User{}
	u := &UnixService{
		logger: logger,
	}

	createUser(t, dto.GetUsername())
	defer deleteUser(t, dto.GetUsername())

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

	createUser(t, dto.GetUsername())

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
