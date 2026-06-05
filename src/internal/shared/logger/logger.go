package logger

import (
	"context"
)

type Logger interface {
	Info(addSource bool, tx context.Context, message string)
	Warn(addSource bool, tx context.Context, message string)
	Error(addSource bool, tx context.Context, message string)
}
