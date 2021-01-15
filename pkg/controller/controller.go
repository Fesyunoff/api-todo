package controller

import (
	"context"

	"github.com/fesyunoff/api/pkg/controller/dto"
)

type Service interface {
	// storage.Todo

	Get(ctx context.Context, task dto.Task) (out []*dto.Task, msg string, err error)
	Add(ctx context.Context, task dto.Task) (msg string, err error)
	Update(ctx context.Context, task dto.Task) (msg string, err error)
	Delete(ctx context.Context, task dto.Task) (msg string, err error)

	CreateUser(ctx context.Context, user dto.User) (msg string, err error)
	GetUsers(ctx context.Context, user dto.User) (out []*dto.User, msg string, err error)
	DeleteUser(ctx context.Context, user dto.User) (msg string, err error)
}
