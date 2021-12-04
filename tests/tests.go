package tests

import (
	"context"
	"os/exec"
	"strings"
	"testing"
)


func DeleteUser(t *testing.T, username string) {
	cmd := exec.CommandContext(
		context.Background(),
		"userdel",
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

func CreateUser(t *testing.T, username string, groups ...string) {
	flags := []string{
		"--create-home",
		"--password",
		"password",
		"--shell",
		"/bin/sh",
		"--user-group",
	}

	if len(groups) > 0 {
		flags = append(flags, "--groups", strings.Join(groups, ","))
	}

	flags = append(flags, username)

	cmd := exec.CommandContext(
		context.Background(),
		"useradd",
		flags...,
	)

	err := cmd.Run()
	if err != nil {
		t.Errorf("Error while creating user: %v\n", err)
		t.FailNow()
	}
}
