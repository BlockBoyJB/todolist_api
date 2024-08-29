package dbmodel

import "time"

type Task struct {
	Id          int       `db:"id"`
	Username    string    `db:"username"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	DueDate     time.Time `db:"due_date"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
