package logger

import (
	"context"
	"fmt"
	stdLog "log"
	"os"
	"strings"
)

type Logger interface {
	Error(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Debug(ctx context.Context, args ...interface{})
}

type log struct {
	l *stdLog.Logger
}

func (l *log) print(ctx context.Context, args ...interface{}) {
	var result strings.Builder

	// TODO Log context

	for _, arg := range append(args) {
		if _, err := result.WriteString(fmt.Sprintf(" %v ", arg)); err != nil {
			panic(err)
		}
	}

	l.l.Println(result.String())
}

func (l *log) Error(ctx context.Context, args ...interface{}) {
	l.print(ctx, args...)
}

func (l *log) Info(ctx context.Context, args ...interface{}) {
	l.print(ctx, args...)
}

func (l *log) Debug(ctx context.Context, args ...interface{}) {
	l.print(ctx, args...)
}

func New(prefix string) Logger {
	return &log{
		l: stdLog.New(os.Stdout, fmt.Sprintf("%s: ", prefix), stdLog.LstdFlags),
	}
}
