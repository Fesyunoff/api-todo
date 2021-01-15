package dto

type Task struct {
	TaskId  int    `json:"task_id,omitempty"`
	Title   string `json:"title"`
	Note    string `json:"note"`
	DueDate string `json:"due_date"`
	UserId  int    `json:"user_id,omitempty"`
}

type User struct {
	UserId int    `json:"user_id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
}
