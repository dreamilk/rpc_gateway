package logger

import (
	"context"
	"fmt"
)

func Info(ctx context.Context, msg string) {
	fmt.Println(msg, "")
}

func Warn(ctx context.Context, msg string) {

}

func Error(ctx context.Context, msg string) {

}
