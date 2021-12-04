package user

type Logger interface {
	Print(string, ...interface{})
}
