package user

import (
	"context"
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"

	mocks "github.com/SSH-Management/linux-user/__mocks__"
)

func deleteUser(username string) {
	cmd := exec.CommandContext(
		context.Background(),
		DeleteCommand,
		"-f",
		"-r",
		username,
	)

	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error while deleting user: %v\n", err)
	}
}

func createUser(username string) {
	cmd := exec.CommandContext(
		context.Background(),
		AddCommand,
		"--create-home",
		"--password",
		"password",
		"--groups",
		"sudo",
		"--shell",
		"/bin/sh",
		"--user-group",
		username,
	)

	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error while creating user: %v\n", err)
	}
}

func TestCreateUser_Success(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	dto := &mocks.User{}
	u := &UnixService{}

	defer deleteUser(dto.GetUsername())

	linuxUser, err := u.Create(context.Background(), dto)

	assert.NoError(err)
	assert.Equal("/home/"+dto.GetUsername(), linuxUser.HomeFolder)
	assert.GreaterOrEqual(linuxUser.UserId, 1000)
}

func TestCreateUser_AlreadyExists(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	strErr := fmt.Sprintf("error while running the create user command, Exit Status: %d, Error: %v",9, createCommandFail[9])

	logger := mocks.NewLogger()
	dto := &mocks.User{}
	u := &UnixService{
		logger: logger,
	}

	createUser(dto.GetUsername())
	defer deleteUser(dto.GetUsername())

	_, err := u.Create(context.Background(), dto)

	assert.Error(err)
	assert.ErrorIs(err, createCommandFail[9])

	assert.Equal(logger.Len(), 1)
	assert.Equal(strErr, logger.Data()[0])
}
