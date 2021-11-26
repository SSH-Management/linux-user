package user

import (
	"context"
	"os/exec"
	osuser "os/user"
	"strconv"
	"strings"
)

var _ UnixInterface = &UnixService{}

type (
	UnixService struct {
		groups map[string]string
		logger Logger
	}

	Logger interface{
		Print(string, ...interface{})
	}

	User struct {
		UserId     int
		GroupId    int
		HomeFolder string
	}

	UnixInterface interface {
		Create(context.Context, Creater) (User, error)
		Delete(context.Context, string) error
	}
)

func NewUnixService(groups map[string]string, logger ...Logger) *UnixService {
	var l Logger = nil

	if len(logger) > 0 {
		l = logger[0]
	}

	return &UnixService{
		groups: groups,
		logger: l,
	}
}

func (s UnixService) mapSystemGroups(groups []string) []string {
	newGroups := make([]string, 0, len(groups))

	for _, group := range groups {
		if replace, ok := s.groups[group]; ok {
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

	if status := cmd.ProcessState.ExitCode(); status != 0 {
		err, ok := createCommandFail[status]

		if ok {
			if s.logger != nil {
				s.logger.Print("error while running the create user command, Exit Status: %d, Error: %v", status, err)
			}

			return User{}, err
		}

		return User{}, ErrCommandFailed
	}

	if err != nil {
		if s.logger != nil {
			s.logger.Print("error while running the create user command: %v", err)
		}
		return User{}, ErrCommandFailed
	}

	u, err := osuser.Lookup(user.GetUsername())

	if err != nil {
		if s.logger != nil {
			s.logger.Print("cannot find user on the system: %v", err)
		}

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

	if status := cmd.ProcessState.ExitCode(); status != 0 {
		err, ok := deleteCommandFail[status]

		if !ok {
			return ErrCommandFailed
		}

		if s.logger != nil {
			s.logger.Print("error while running the delete user command, Exit Status: %d, Error: %v", status, err)
		}

		return err
	}

	if err != nil {
		if s.logger != nil {
			s.logger.Print("error while running the delete user command: %v", err)
		}

		return ErrCommandFailed
	}

	return nil
}
