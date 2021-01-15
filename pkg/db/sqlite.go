package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/fesyunoff/api/pkg/controller/dto"
	"github.com/fesyunoff/api/pkg/storage"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteTodoStorage struct {
	Conn *sql.DB
}

var _ storage.Todo = (*SQLiteTodoStorage)(nil)

func CreateSQLiteDB(name string) (sqliteDB *sql.DB) {
	path := fmt.Sprintf("./%s", name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("Creating %s...", name)
		file, err := os.Create(name)
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Printf("%s created", name)
	} else {
		log.Printf("%s exist", name)
	}

	sqliteDB, _ = sql.Open("sqlite3", path)

	// do not forget close connection at func call place
	// defer sqliteDB.Close()

	createTables(sqliteDB)

	return
}

//Create new entry in SQLite database table "tasks" with data from "t"
func (s *SQLiteTodoStorage) CreateTask(db *sql.DB, u dto.User, t dto.Task) (msg string, err error) {
	req := `INSERT INTO tasks(title, note, due_date, user_id) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	result, err := statement.Exec(t.Title, t.Note, t.DueDate, u.UserId)
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		raws, _ := result.RowsAffected()
		taskId, _ := result.LastInsertId()
		msg = fmt.Sprintf("create task with id %d; %d raws affected", taskId, raws)
	}
	return
}

//Update entry with "task_id" == "t.TaskId' in SQLite database table "tasks"
func (s *SQLiteTodoStorage) UpdateTask(db *sql.DB, u dto.User, t dto.Task) (msg string, err error) {
	var req string
	if u.Role == "adm" {
		req = `UPDATE tasks SET title = ?, note = ?, due_date = ?, user_id = ? WHERE task_id = ?;`
	} else {
		req = `UPDATE tasks SET title = ?, note = ?, due_date = ? WHERE user_id = ? AND task_id = ?;`
	}
	statement, err := db.Prepare(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	result, err := statement.Exec(t.Title, t.Note, t.DueDate, u.UserId, t.TaskId)
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		raws, _ := result.RowsAffected()
		taskId, _ := result.LastInsertId()
		msg = fmt.Sprintf("task %d update, %d raws affected", taskId, raws)
	}
	return
}

//Delete entry with "task_id" == "t.TaskId' in SQLite database table "tasks"
func (s *SQLiteTodoStorage) DeleteTask(db *sql.DB, u dto.User, t dto.Task) (msg string, err error) {
	var req string
	if u.Role == "adm" {
		req = `DELETE FROM tasks WHERE task_id = ?;`
	} else {
		req = fmt.Sprintf("DELETE FROM tasks WHERE user_id = %d AND task_id = ?;", u.UserId)
	}
	statement, err := db.Prepare(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	result, err := statement.Exec(t.TaskId)
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		raws, _ := result.RowsAffected()
		msg = fmt.Sprintf("%d raws affected", raws)
	}
	return
}

//Return entry with "task_id" == "t.TaskId'(all if "0") from SQLite database table "tasks"
func (s *SQLiteTodoStorage) DisplayTask(db *sql.DB, u dto.User, t dto.Task) (out []*dto.Task, err error) {
	var req string
	if u.Role == "adm" {
		if t.TaskId == 0 {
			req = "SELECT * FROM tasks ORDER BY task_id;"
		} else {
			req = fmt.Sprintf("SELECT * FROM tasks WHERE task_id = %d;", t.TaskId)
		}
	} else {
		if t.TaskId == 0 {
			req = fmt.Sprintf("SELECT * FROM tasks WHERE user_id = %d ORDER BY task_id;", u.UserId)
		} else {
			req = fmt.Sprintf("SELECT * FROM tasks WHERE user_id = %d AND task_id = %d;", u.UserId, t.TaskId)
		}
	}
	row, err := db.Query(req)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		a := dto.Task{}
		err = row.Scan(&a.TaskId, &a.Title, &a.Note, &a.DueDate, &a.UserId)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(a)
		out = append(out, &a)
	}
	return
}

func (s *SQLiteTodoStorage) CreateUser(db *sql.DB, u dto.User) (msg string, err error) {
	req := `INSERT INTO users(name, role) VALUES (?, ?)`
	statement, err := db.Prepare(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if u.Role == "usr" || u.Role == "adm" {
		result, err := statement.Exec(u.Name, u.Role)
		if err != nil {
			log.Fatalln(err.Error())
		}
		id, _ := result.LastInsertId()
		raws, err := result.RowsAffected()
		if err == nil {
			msg = fmt.Sprintf("create user with id %d; %d raws affected", id, raws)
		}
	} else {
		msg = fmt.Sprintln("ERROR: uncorrect format: use 'adm' or 'usr' to field 'role'")
	}
	return
}
func (s *SQLiteTodoStorage) DisplayUsers(db *sql.DB) (out []*dto.User, err error) {
	req := "SELECT * FROM users ORDER BY user_id;"
	row, err := db.Query(req)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		u := dto.User{}
		err = row.Scan(&u.UserId, &u.Name, &u.Role)
		if err != nil {
			log.Fatal(err)
		}
		out = append(out, &u)
	}
	return
}

func (s *SQLiteTodoStorage) ReturnUser(db *sql.DB, id int64) (out dto.User, err error) {
	req := fmt.Sprintf("SELECT * FROM users WHERE user_id = %d;", id)
	// fmt.Println(req)
	row := db.QueryRow(req)
	usr := dto.User{}
	err = row.Scan(&usr.UserId, &usr.Name, &usr.Role)
	if err != nil {
		log.Println(err.Error())
	}
	out = usr
	return
}

func (s *SQLiteTodoStorage) DeleteUser(db *sql.DB, u dto.User) (msg string, err error) {
	req := `DELETE FROM users WHERE user_id = ?;`
	statement, err := db.Prepare(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	result, err := statement.Exec(u.UserId)
	if err != nil {
		log.Fatalln(err.Error())
	}
	raws, err := result.RowsAffected()
	if err == nil {
		msg = fmt.Sprintf("%d raws affected", raws)
	}
	return
}

func createTables(db *sql.DB) {
	createTaskTableReq := `CREATE TABLE IF NOT EXISTS tasks (
		"task_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"title" TEXT,
		"note" TEXT,
		"due_date" TEXT,
		"user_id" INTEGER
	  );`
	createNewTable(db, "tasks", createTaskTableReq)

	createUserTableReq := `CREATE TABLE IF NOT EXISTS users (
		"user_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"name" TEXT,
		"role" TEXT
	  );`
	createNewTable(db, "users", createUserTableReq)
}

func createNewTable(db *sql.DB, name string, createReq string) {
	existReq := fmt.Sprintf("SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = '%s';", name)
	row, err := db.Query(existReq)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var a int
	for row.Next() {
		err = row.Scan(&a)
		if err != nil {
			log.Fatal(err)
		}
	}
	if a == 0 {
		log.Printf("Create table %s...", name)
		statement, err := db.Prepare(createReq)
		if err != nil {
			log.Fatal(err.Error())
		}
		_, err = statement.Exec()
		if err != nil {
			log.Fatalln(err.Error())
			log.Printf("ERROR: table %s not created", name)
		} else {
			log.Printf("%s created", name)
		}

		if name == "users" {
			insertRootUser(db)
		}

	} else {
		log.Printf("table '%s' exist", name)
	}
}

func insertRootUser(db *sql.DB) {
	statement, err := db.Prepare(`INSERT INTO users(name, role) VALUES ("root", "adm");`)
	if err != nil {
		log.Fatal(err.Error())
	}
	resp, err := statement.Exec()
	if err != nil {
		log.Fatalln(err.Error())
		log.Printf("ERROR: user 'root' not created")
	} else {
		id, _ := resp.LastInsertId()
		log.Printf("user 'root'  with id %d created", id)
	}
}

func NewSQLiteTodoStorage(conn *sql.DB) *SQLiteTodoStorage {
	return &SQLiteTodoStorage{
		Conn: conn,
	}
}
