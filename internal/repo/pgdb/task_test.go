package pgdb

import (
	"github.com/jackc/pgx/v5"
	"time"
	"todolist_api/internal/model/dbmodel"
	"todolist_api/internal/repo/pgerrs"
)

func (s *pgdbTestSuite) setupTestsData() string {
	username := "vasya"
	sql, args, _ := s.pg.Builder.
		Insert("\"user\"").
		Columns("username", "password").
		Values(username, "abc").
		ToSql()

	if _, err := s.pg.Pool.Exec(s.ctx, sql, args...); err != nil {
		panic(err)
	}
	return username
}

func (s *pgdbTestSuite) TestTaskRepo_Create() {
	username := s.setupTestsData()

	testCases := []struct {
		testName  string
		task      *dbmodel.Task
		expectErr error
	}{
		{
			testName: "Correct test",
			task: &dbmodel.Task{
				Username:    username,
				Title:       "Test",
				Description: "desc test",
				DueDate:     time.Date(2024, 8, 27, 23, 0, 0, 0, time.UTC),
			},
			expectErr: nil,
		},
		{
			testName: "User not exist",
			task: &dbmodel.Task{
				Username:    "foobar",
				Title:       "title",
				Description: "desc",
				DueDate:     time.Now(),
			},
			expectErr: pgerrs.ErrForeignKey,
		},
	}

	for _, tc := range testCases {
		err := s.task.Create(s.ctx, tc.task)
		s.Assert().Equal(tc.expectErr, err)

		if tc.expectErr == nil {
			sql, args, _ := s.pg.Builder.
				Select("*").
				From("task").
				Where("id = ?", tc.task.Id).
				ToSql()

			var actualTask dbmodel.Task
			err = s.pg.Pool.QueryRow(s.ctx, sql, args...).Scan(
				&actualTask.Id,
				&actualTask.Username,
				&actualTask.Title,
				&actualTask.Description,
				&actualTask.DueDate,
				&actualTask.CreatedAt,
				&actualTask.UpdatedAt,
			)
			s.Assert().Nil(err)
			s.Assert().Equal(*tc.task, actualTask)
		}
	}
}

func (s *pgdbTestSuite) TestTaskRepo_Update() {
	username := s.setupTestsData()
	task := &dbmodel.Task{

		Username:    username,
		Title:       "Hello",
		Description: "desc hello",
		DueDate:     time.Date(2024, 8, 27, 23, 0, 0, 0, time.UTC),
	}
	if err := s.task.Create(s.ctx, task); err != nil { // так можно делать, потому что есть отдельный тест
		panic(err)
	}

	testCases := []struct {
		testName  string
		task      *dbmodel.Task
		expectErr error
	}{
		{
			testName: "Correct test",
			task: &dbmodel.Task{
				Id:          task.Id,
				Username:    task.Username,
				Title:       "New title",
				Description: "New desc",
				DueDate:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			expectErr: nil,
		},
		{
			testName: "User not exist",
			task: &dbmodel.Task{
				Id:          task.Id,
				Username:    "foobar",
				Title:       "New title",
				Description: "New desc",
				DueDate:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectErr: pgerrs.ErrNotFound,
		},
		{
			testName: "Id not exist",
			task: &dbmodel.Task{
				Id:          123123123,
				Username:    username,
				Title:       "New title",
				Description: "New desc",
				DueDate:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectErr: pgerrs.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		err := s.task.Update(s.ctx, tc.task)
		s.Assert().Equal(tc.expectErr, err)

		if tc.expectErr == nil {
			sql, args, _ := s.pg.Builder.
				Select("id", "username", "title", "description", "due_date", "created_at", "updated_at").
				From("task").
				Where("id = ?", tc.task.Id).
				Where("username = ?", tc.task.Username).
				ToSql()

			var actualTask dbmodel.Task
			err = s.pg.Pool.QueryRow(s.ctx, sql, args...).Scan(
				&actualTask.Id,
				&actualTask.Username,
				&actualTask.Title,
				&actualTask.Description,
				&actualTask.DueDate,
				&actualTask.CreatedAt,
				&actualTask.UpdatedAt,
			)
			s.Assert().Nil(err)
			s.Assert().Equal(*tc.task, actualTask)
			s.Assert().NotEqual(task.UpdatedAt, actualTask.UpdatedAt)
		}
	}
}

func (s *pgdbTestSuite) TestTaskRepo_Delete() {
	username := s.setupTestsData()
	task := &dbmodel.Task{

		Username:    username,
		Title:       "Hello",
		Description: "desc hello",
		DueDate:     time.Date(2024, 8, 27, 23, 0, 0, 0, time.UTC),
	}
	if err := s.task.Create(s.ctx, task); err != nil { // так можно делать, потому что есть отдельный тест
		panic(err)
	}

	testCases := []struct {
		testName  string
		id        int
		username  string
		expectErr error
	}{
		{
			testName:  "User not exist",
			id:        task.Id,
			username:  "petya",
			expectErr: pgerrs.ErrNotFound,
		},
		{
			testName:  "Id not exist",
			id:        12312312,
			username:  username,
			expectErr: pgerrs.ErrNotFound,
		},
		{
			testName:  "Correct test",
			id:        task.Id,
			username:  username,
			expectErr: nil,
		},
		{
			testName:  "Task has been deleted",
			id:        task.Id,
			username:  username,
			expectErr: pgerrs.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		err := s.task.Delete(s.ctx, tc.id, tc.username)
		s.Assert().Equal(tc.expectErr, err)

		if tc.expectErr == nil {
			sql, args, _ := s.pg.Builder.
				Select("*").
				From("task").
				Where("id = ?", tc.id).
				ToSql()

			err = s.pg.Pool.QueryRow(s.ctx, sql, args...).Scan()
			s.Assert().Equal(pgx.ErrNoRows, err)
		}
	}
}
