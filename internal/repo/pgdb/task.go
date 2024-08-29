package pgdb

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"todolist_api/internal/model/dbmodel"
	"todolist_api/internal/repo/pgerrs"
	"todolist_api/pkg/postgres"
)

type TaskRepo struct {
	*postgres.Postgres
}

func NewTaskRepo(pg *postgres.Postgres) *TaskRepo {
	return &TaskRepo{pg}
}

// Create создает запись.
// На вход принимает task с полями: Username, Title, Description, DueDate
// Остальные поля структуры обновляются после записи в бд (поэтому в аргументах пойнтер)
func (r *TaskRepo) Create(ctx context.Context, t *dbmodel.Task) error {
	sql, args, _ := r.Builder.
		Insert("task").
		Columns("username", "title", "description", "due_date").
		Values(t.Username, t.Title, t.Description, t.DueDate).
		Suffix("returning id, created_at, updated_at").
		ToSql()

	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&t.Id,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23503" {
				return pgerrs.ErrForeignKey
			}
		}
		return err
	}
	return nil
}

func (r *TaskRepo) Find(ctx context.Context, username string) ([]dbmodel.Task, error) {
	sql, args, _ := r.Builder.
		Select("id", "title", "description", "due_date", "created_at", "updated_at").
		From("task").
		Where("username = ?", username).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dbmodel.Task
	for rows.Next() {
		var task dbmodel.Task

		err = rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, task)
	}
	return result, nil
}

func (r *TaskRepo) FindById(ctx context.Context, id int, username string) (dbmodel.Task, error) {
	sql, args, _ := r.Builder.
		Select("title", "description", "due_date", "created_at", "updated_at").
		From("task").
		Where("id = ?", id).
		Where("username = ?", username).
		ToSql()

	var task dbmodel.Task
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dbmodel.Task{}, pgerrs.ErrNotFound
		}
		return dbmodel.Task{}, err
	}
	return task, nil
}

// Update обновляет запись (по-сути, все доступные поля заменяет на новые).
// На вход принимает task с полями: Id, Username, Title, Description, DueDate.
// Поля CreatedAt, UpdatedAt устанавливаются в соответствии с бд аналогично функции Create (через пойнтер)
func (r *TaskRepo) Update(ctx context.Context, t *dbmodel.Task) error {
	sql, args, _ := r.Builder.
		Update("task").
		Set("title", t.Title).
		Set("description", t.Description).
		Set("due_date", t.DueDate).
		Where("id = ?", t.Id).
		Where("username = ?", t.Username).
		Suffix("returning created_at, updated_at").
		ToSql()

	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgerrs.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *TaskRepo) Delete(ctx context.Context, id int, username string) error {
	sql, args, _ := r.Builder.
		Delete("task").
		Where("id = ?", id).
		Where("username = ?", username).
		Suffix("returning id").
		ToSql()

	// так делаем, потому что по заданию надо возвращать ошибку если задачи не было,
	// а если делать простой Exec, то он никакую ошибку не вернет
	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(nil); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgerrs.ErrNotFound
		}
		return err
	}
	return nil
}
