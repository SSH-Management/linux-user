package mocks

import (
	"fmt"
	"math/rand"
)

type User struct{
	username string
}

func (u User) GetSystemGroups() []string {
	return []string{"sudo"}
}

func (u User) GetPlainTextPassword() string {
	return "password"
}

func (u User) GetDefaultShell() string {
	return "/bin/sh"
}

func (u *User) GetUsername() string {
	if u.username != "" {
		return u.username
	}

	u.username = fmt.Sprintf("test_user_%d", rand.Int31n(10_000))

	return u.username
}
