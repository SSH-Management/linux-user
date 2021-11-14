package user

import "errors"

var (
	ErrCommandFailed     = errors.New("command failed")
	ErrUserAlreadyExists = errors.New("user already exists")

	createCommandFail = map[int]error{
		1:  errors.New("can't update password file"),
		2:  errors.New("invalid command syntax"),
		3:  errors.New("invalid argument to option"),
		4:  errors.New("UID already in use (and no -o)"),
		6:  errors.New("specified group doesn't exist"),
		9:  ErrUserAlreadyExists,
		10: errors.New("can't update group file"),
		12: errors.New("can't create home directory"),
	}

	deleteCommandFail = map[int]error{
		1:  errors.New("can't update password file"),
		2:  errors.New("invalid command syntax"),
		6:  errors.New("specified user doesn't exist"),
		8:  errors.New("user currently logged in"),
		10: errors.New("can't update group file"),
		12: errors.New("can't create home directory"),
	}
)
