package repo

import (
	"context"
	"todolist_api/internal/model/dbmodel"
	"todolist_api/internal/repo/pgdb"
	"todolist_api/pkg/postgres"
)

type User interface {
	Create(ctx context.Context, u dbmodel.User) error
	FindByUsername(ctx context.Context, username string) (dbmodel.User, error)
}

type Task interface {
	Create(ctx context.Context, t *dbmodel.Task) error
	Find(ctx context.Context, username string) ([]dbmodel.Task, error)
	FindById(ctx context.Context, id int, username string) (dbmodel.Task, error)
	Update(ctx context.Context, t *dbmodel.Task) error
	Delete(ctx context.Context, id int, username string) error
}

type Repositories struct {
	User
	Task
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User: pgdb.NewUserRepo(pg),
		Task: pgdb.NewTaskRepo(pg),
	}
}
