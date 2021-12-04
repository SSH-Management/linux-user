package user

import (
	"context"
	"testing"

	mocks "github.com/SSH-Management/linux-user/__mocks__"
	"github.com/SSH-Management/linux-user/tests"
	"github.com/stretchr/testify/require"
)


func TestGetUserGroups(t *testing.T) {
	t.Parallel()

	dto := &mocks.User{}
	assert := require.New(t)

	t.Cleanup(func() {
		tests.DeleteUser(t, dto.GetUsername())
	})

	tests.CreateUser(t, dto.GetUsername(), "sudo")

	u := User{
		Username: dto.GetUsername(),
	}

	groups, err := u.GetSystemGroups(context.Background())

	assert.NoError(err)
	assert.NotEmpty(groups)
	assert.Contains(groups, "sudo")
	assert.Contains(groups, dto.GetUsername())
}

func TestGetUserGroups_MappedGroups(t *testing.T) {
	t.Parallel()

	dto := &mocks.User{}
	assert := require.New(t)

	t.Cleanup(func() {
		tests.DeleteUser(t, dto.GetUsername())
	})

	tests.CreateUser(t, dto.GetUsername(), "sudo")

	u := User{
		Username: dto.GetUsername(),
		invertedGroups: map[string]string{
			"sudo": "my_sudo_map",
		},
	}

	groups, err := u.GetSystemGroups(context.Background())

	assert.NoError(err)
	assert.NotEmpty(groups)
	assert.NotContains(groups, "sudo")
	assert.Contains(groups, dto.GetUsername())
	assert.Contains(groups, "my_sudo_map")
}


func TestGetUserGroups_UserNotFound(t *testing.T) {
	t.Parallel()

	dto := &mocks.User{}
	assert := require.New(t)

	u := User{
		Username: dto.GetUsername(),
	}

	groups, err := u.GetSystemGroups(context.Background())

	assert.Nil(groups)
	assert.Error(err)
	assert.ErrorIs(err, ErrUserDoesNotExist)
}


func TestGetUserGroups_UsernameEmpty(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	u := User{
		Username: "",
	}

	groups, err := u.GetSystemGroups(context.Background())

	assert.Nil(groups)
	assert.Error(err)
	assert.ErrorIs(err, ErrUsernameEmpty)
}


