package user

type Creater interface {
	GetSystemGroups() []string
	GetPlainTextPassword() string
	GetDefaultShell() string
	GetUsername() string
}
