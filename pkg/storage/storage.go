package storage

import (
	"database/sql"

	"github.com/fesyunoff/api/pkg/controller/dto"
)

type Todo interface {
	// Create(db *sql.DB, t dto.Task) (int64, error)
	CreateTask(db *sql.DB, u dto.User, t dto.Task) (msg string, err error)
	DisplayTask(db *sql.DB, u dto.User, t dto.Task) (out []*dto.Task, err error)
	UpdateTask(db *sql.DB, u dto.User, t dto.Task) (msg string, err error)
	DeleteTask(db *sql.DB, u dto.User, t dto.Task) (msg string, err error)

	CreateUser(db *sql.DB, u dto.User) (msg string, err error)
	DisplayUsers(db *sql.DB) (out []*dto.User, err error)
	ReturnUser(db *sql.DB, id int64) (u dto.User, err error)
	DeleteUser(db *sql.DB, u dto.User) (msg string, err error)
}
