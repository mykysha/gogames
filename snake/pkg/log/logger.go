package log

type Logger interface {
	Error(msg string, args ...any)
}
