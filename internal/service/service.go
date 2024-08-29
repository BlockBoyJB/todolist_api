package service

import (
	"context"
	"time"
	"todolist_api/internal/repo"
	"todolist_api/pkg/hasher"
)

type (
	UserInput struct {
		Username string
		Password string
	}
	TaskCreateInput struct {
		Username    string
		Title       string
		Description string
		DueDate     time.Time
	}
	TaskUpdateInput struct {
		Id          int
		Username    string
		Title       string
		Description string
		DueDate     time.Time
	}
	TaskOutput struct {
		Id          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}
)

type Auth interface {
	CreateToken(username string) (string, error)
	ParseToken(tokenString string) (*TokenClaims, error)
}

type User interface {
	Create(ctx context.Context, input UserInput) error
	VerifyPassword(ctx context.Context, input UserInput) (bool, error)
}

type Task interface {
	Create(ctx context.Context, input TaskCreateInput) (TaskOutput, error)
	Find(ctx context.Context, username string) ([]TaskOutput, error)
	FindById(ctx context.Context, id int, username string) (TaskOutput, error)
	Update(ctx context.Context, input TaskUpdateInput) (TaskOutput, error)
	Delete(ctx context.Context, id int, username string) error
}

type (
	Services struct {
		Auth Auth
		User User
		Task Task
	}
	ServicesDependencies struct {
		Repos    *repo.Repositories
		Hasher   hasher.Hasher
		SignKey  string
		TokenTTL time.Duration
	}
)

func NewServices(d *ServicesDependencies) *Services {
	return &Services{
		Auth: newAuthService(d.SignKey, d.TokenTTL),
		User: newUserService(d.Repos.User, d.Hasher),
		Task: newTaskService(d.Repos.Task),
	}
}
