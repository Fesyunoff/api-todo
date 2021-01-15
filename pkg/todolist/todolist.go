package todolist

import (
	"context"
	"strconv"

	"github.com/fesyunoff/api/pkg/auth"
	"github.com/fesyunoff/api/pkg/controller"
	"github.com/fesyunoff/api/pkg/controller/dto"
	"github.com/fesyunoff/api/pkg/db"
	"github.com/fesyunoff/api/pkg/msg"
)

type Service struct {
	db  *db.SQLiteTodoStorage
	msg *msg.Messanger
}

var _ controller.Service = (*Service)(nil)

func (s *Service) Add(ctx context.Context, p dto.Task) (msg string, err error) {

	_, msg, err = executeTaskMtd("add", ctx, p, s)

	return
}

func (s *Service) Get(ctx context.Context, p dto.Task) (out []*dto.Task, msg string, err error) {

	out, msg, err = executeTaskMtd("get", ctx, p, s)

	return
}

func (s *Service) Update(ctx context.Context, p dto.Task) (msg string, err error) {

	_, msg, err = executeTaskMtd("update", ctx, p, s)

	return
}

func (s *Service) Delete(ctx context.Context, p dto.Task) (msg string, err error) {

	_, msg, err = executeTaskMtd("delete", ctx, p, s)

	return
}

func (s *Service) CreateUser(ctx context.Context, u dto.User) (msg string, err error) {

	_, msg, err = executeUserMtd("createUser", ctx, u, s)

	return
}

func (s *Service) GetUsers(ctx context.Context, u dto.User) (out []*dto.User, msg string, err error) {

	out, msg, err = executeUserMtd("getUsers", ctx, u, s)

	return
}

func (s *Service) DeleteUser(ctx context.Context, u dto.User) (msg string, err error) {

	_, msg, err = executeUserMtd("deleteUser", ctx, u, s)

	return
}

func executeTaskMtd(mtd string, ctx context.Context, p dto.Task, s *Service) (out []*dto.Task, msg string, err error) {
	idStr := auth.PersonIDFromContext(ctx)
	id, _ := strconv.ParseInt(idStr, 10, 64)
	usr, err := s.db.ReturnUser(s.db.Conn, id)
	if err != nil || usr.UserId == 0 {
		msg = "ERROR: permission denied"
		return
	} else {
		switch mtd {
		case "add":
			msg, err = s.db.CreateTask(s.db.Conn, usr, p)
		case "get":
			out, err = s.db.DisplayTask(s.db.Conn, usr, p)
		case "update":
			msg, err = s.db.UpdateTask(s.db.Conn, usr, p)
		case "delete":
			msg, err = s.db.DeleteTask(s.db.Conn, usr, p)
		}
	}
	s.msg.SentMessage(msg)
	return
}

func executeUserMtd(mtd string, ctx context.Context, u dto.User, s *Service) (out []*dto.User, msg string, err error) {
	idStr := auth.PersonIDFromContext(ctx)
	id, _ := strconv.ParseInt(idStr, 10, 64)
	// fmt.Println(id)
	usr, err := s.db.ReturnUser(s.db.Conn, id)
	// fmt.Println(usr)
	if err != nil || usr.Role != "adm" || usr.UserId == 0 {
		msg = "ERROR: permission denied"
	} else {
		switch mtd {
		case "createUser":
			msg, err = s.db.CreateUser(s.db.Conn, u)
		case "getUsers":
			out, err = s.db.DisplayUsers(s.db.Conn)
		case "deleteUser":
			msg, err = s.db.DeleteUser(s.db.Conn, u)
		}
	}
	s.msg.SentMessage(msg)
	return
}

func NewService(db *db.SQLiteTodoStorage, msg *msg.Messanger) *Service {
	return &Service{
		db:  db,
		msg: msg,
	}
}
