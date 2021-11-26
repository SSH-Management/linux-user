package user

import "errors"

var (
	ErrCommandFailed         = errors.New("command failed")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrCantUpdateGroupFile   = errors.New("can't update group file")
	ErrCantUpdatePasswdFile  = errors.New("can't update password file")
	ErrInvalidCommandSytax   = errors.New("invalid command syntax")
	ErrCantDeleteHomeDir     = errors.New("can't delete home directory")
	ErrCantCreateHomeDir     = errors.New("can't create home directory")
	ErrUserIdAlreadyInUse    = errors.New("UID already in use (and no -o)")
	ErrUserLoggedIn          = errors.New("user currently logged in")
	ErrUserDoesNotExist      = errors.New("specified user doesn't exist")
	ErrGroupDoesNotExist     = errors.New("specified group doesn't exist")
	ErrInvalidArgumentOption = errors.New("invalid argument to option")

	createCommandFail = map[int]error{
		1:  ErrCantUpdatePasswdFile,
		2:  ErrInvalidCommandSytax,
		3:  ErrInvalidArgumentOption,
		4:  ErrUserIdAlreadyInUse,
		6:  ErrGroupDoesNotExist,
		9:  ErrUserAlreadyExists,
		10: ErrCantUpdateGroupFile,
		12: ErrCantCreateHomeDir,
	}

	deleteCommandFail = map[int]error{
		1:  ErrCantUpdatePasswdFile,
		2:  ErrInvalidCommandSytax,
		6:  ErrUserDoesNotExist,
		8:  ErrUserLoggedIn,
		10: ErrCantUpdateGroupFile,
		12: ErrCantDeleteHomeDir,
	}
)
