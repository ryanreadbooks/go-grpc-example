package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 检查context.Context是否有错误
func CheckContext(ctx context.Context) (err error) {
	if errors.Is(ctx.Err(), context.Canceled) {
		err = status.Error(codes.Canceled, "canceled")
		return
	}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		err = status.Error(codes.DeadlineExceeded, "deadline exceeded")
		return
	}
	err = ctx.Err()

	return
}

func CheckUUIDValid(id string) error {
	_, err := uuid.Parse(id)
	return err
}
