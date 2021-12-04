package user

import (
	"context"
	"errors"
	"os/user"
)

var ErrUsernameEmpty = errors.New("username cannot be empty string")

type User struct {
	UserId     int
	GroupId    int
	HomeFolder string
	Username   string

	groups         []string
	invertedGroups map[string]string
}

func (s User) GetSystemGroups(ctx context.Context) ([]string, error) {
	if s.groups != nil {
		return s.groups, nil
	}

	if s.Username == "" {
		return nil, ErrUsernameEmpty
	}

	u, err := user.Lookup(s.Username)
	if err != nil {
		if _, ok := err.(user.UnknownUserError); ok {
			return nil, ErrUserDoesNotExist
		}
		return nil, ErrCommandFailed
	}

	groups, err := u.GroupIds()
	if err != nil {
		return nil, ErrGroupDoesNotExist
	}

	systemGroups := make([]string, 0, len(groups))

	for _, id := range groups {
		group, err := user.LookupGroupId(id)
		if err != nil {
			return nil, err
		}

		systemGroups = append(systemGroups, group.Name)
	}

	s.groups = mapSystemGroups(systemGroups, s.invertedGroups)

	return s.groups, nil
}
