package db

import (
	"database/sql"

	"github.com/fesyunoff/api/pkg/controller/dto"
	_ "github.com/go-sql-driver/mysql"
)

/*
CREATE DATABASE tododb;
USE tododb;
CREATE TABLE tasks (
    task_id int auto_increment primary key,
    title varchar(20) not null,
    note varchar(60) not null,
    due_date varchar(10) not null,
    user_id int not null
	);
CREATE TABLE users (
    user_id int auto_increment primary key,
    name varchar(30) not null,
    password varchar(20) not null,
    role varchar(5) not null
*/

type MySQL struct {
	Conn *sql.DB
}

// var _ storage.Todo = (*MySQL)(nil)

func (m *MySQL) CreateConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:pass/mysql1")

	if err != nil {
		panic(err)
	}
	return db, err
}

// defer db.Close()

func (m *MySQL) Create(db *sql.DB, t dto.Task) (int64, error) {
	result, err := db.Exec("INSERT INTO tododb.Tasks (title, note, due_date, user_id) VALUES (?, ?, ?, ?)",
		// "Do", "do enything", "02.02.2021", 1)
		t.Title, t.Note, t.DueDate, t.UserId)
	if err != nil {
		panic(err)
	}
	return result.LastInsertId()
}

func NewMySQL(conn *sql.DB) *MySQL {
	return &MySQL{
		Conn: conn,
	}
}
