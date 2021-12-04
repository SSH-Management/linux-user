package user

import (
	"context"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

var _ UnixInterface = &UnixService{}

type (
	UnixInterface interface {
		Find(context.Context, string) (User, error)
		Create(context.Context, Creater) (User, error)
		Delete(context.Context, string) error
	}

	UnixService struct {
		logger         Logger
		groups         map[string]string
		invertedGroups map[string]string
	}
)

func NewUnixService(groups map[string]string, logger ...Logger) *UnixService {
	var l Logger = nil

	if len(logger) > 0 {
		l = logger[0]
	}

	invertedGroups := make(map[string]string, len(groups))

	for key, value := range groups {
		invertedGroups[value] = key
	}

	return &UnixService{
		groups:         groups,
		invertedGroups: invertedGroups,
		logger:         l,
	}
}

func (s UnixService) Find(ctx context.Context, username string) (User, error) {
	u, err := user.Lookup(username)

	if err != nil {
		if s.logger != nil {
			s.logger.Print("cannot find user on the system: %v", err)
		}

		return User{}, ErrUserDoesNotExist
	}

	userId, _ := strconv.ParseInt(u.Uid, 10, 64)
	groupId, _ := strconv.ParseInt(u.Gid, 10, 64)
	
	return User{
		UserId:     int(userId),
		GroupId:    int(groupId),
		HomeFolder: u.HomeDir,
		Username:   username,

		invertedGroups: s.invertedGroups,
	}, nil
}

func (s UnixService) Create(ctx context.Context, dto Creater) (User, error) {
	systemGroups := mapSystemGroups(dto.GetSystemGroups(), s.groups)

	cmd := exec.CommandContext(
		ctx,
		AddCommand,
		"--create-home",
		"--password",
		dto.GetPlainTextPassword(),
		"--groups",
		strings.Join(systemGroups, ","),
		"--shell",
		dto.GetDefaultShell(),
		"--user-group",
		dto.GetUsername(),
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

	return s.Find(ctx, dto.GetUsername())
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
