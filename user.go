package user

import (
	"context"
	"log"
	"os/exec"
	osuser "os/user"
	"strconv"
	"strings"
)


type (
	UnixService struct {
		sudoGroup map[string]string
	}

	User struct {
		UserId     int
		GroupId    int
		HomeFolder string
	}

	UnixInterface interface {
		Create(context.Context, User) (User, error)
		Delete(context.Context, string) error
	}
)

func NewUnixService(sudoGroup map[string]string, logger *log.Logger) UnixService {
	return UnixService{
		sudoGroup: sudoGroup,
	}
}

func (s UnixService) mapSystemGroups(groups []string) []string {
	newGroups := make([]string, 0, len(groups))

	for _, group := range groups {
		if replace, ok := s.sudoGroup[group]; ok {
			newGroups = append(newGroups, replace)
		} else {
			newGroups = append(newGroups, group)
		}
	}

	return newGroups
}

func (s UnixService) Create(ctx context.Context, user Creater) (User, error) {
	systemGroups := s.mapSystemGroups(user.GetSystemGroups())

	cmd := exec.CommandContext(
		ctx,
		AddCommand,
		"--create-home",
		"--password",
		user.GetPlainTextPassword(),
		"--groups",
		strings.Join(systemGroups, ","),
		"--shell",
		user.GetDefaultShell(),
		"--user-group",
		user.GetUsername(),
	)

	err := cmd.Run()

	if err != nil {
		return User{}, ErrCommandFailed
	}

	if status := cmd.ProcessState.ExitCode(); status != 0 {
		return User{}, createCommandFail[status]
	}

	u, err := osuser.Lookup(user.GetUsername())

	if err != nil {
		return User{},err
	}

	userId, _ := strconv.ParseInt(u.Uid, 10, 64)
	groupId, _ := strconv.ParseInt(u.Gid, 10, 64)


	return User{
		UserId:     int(userId),
		GroupId:    int(groupId),
		HomeFolder: u.HomeDir,
	}, nil
}

func (s UnixService) Delete(ctx context.Context, userName string) error {
	cmd := exec.CommandContext(
		ctx,
		DeleteCommand,
		"-f",
		"-r",
		userName,
	)

	err := cmd.Run()

	if err != nil {
		return err
	}

	if status := cmd.ProcessState.ExitCode(); status != 0 {
		return deleteCommandFail[status]
	}

	return nil
}
