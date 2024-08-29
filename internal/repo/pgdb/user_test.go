package pgdb

import (
	"todolist_api/internal/model/dbmodel"
	"todolist_api/internal/repo/pgerrs"
)

func (s *pgdbTestSuite) TestUserRepo_Create() {
	testCases := []struct {
		testName  string
		user      dbmodel.User
		expectErr error
	}{
		{
			testName: "Correct test",
			user: dbmodel.User{
				Username: "vasya",
				Password: "abc",
			},
			expectErr: nil,
		},
		{
			testName: "User already exist",
			user: dbmodel.User{
				Username: "vasya",
				Password: "foobar",
			},
			expectErr: pgerrs.ErrAlreadyExists,
		},
	}

	for _, tc := range testCases {
		err := s.user.Create(s.ctx, tc.user)
		s.Assert().Equal(tc.expectErr, err)

		if tc.expectErr == nil {
			sql, args, _ := s.pg.Builder.
				Select("username", "password").
				From("\"user\"").
				Where("username = ?", tc.user.Username).
				ToSql()

			var actualUser dbmodel.User

			err = s.pg.Pool.QueryRow(s.ctx, sql, args...).Scan(
				&actualUser.Username,
				&actualUser.Password,
			)
			s.Assert().Nil(err)
			s.Assert().Equal(tc.user, actualUser)
		}
	}
}

func (s *pgdbTestSuite) TestUserRepo_FindByUsername() {
	username, password := "vasya", "abc"
	if err := s.user.Create(s.ctx, dbmodel.User{
		Username: username,
		Password: password,
	}); err != nil {
		panic(err)
	}

	testCases := []struct {
		testName  string
		username  string
		password  string
		expectErr error
	}{
		{
			testName:  "correct test",
			username:  username,
			password:  password,
			expectErr: nil,
		},
		{
			testName:  "user not exist",
			username:  "petya",
			expectErr: pgerrs.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		u, err := s.user.FindByUsername(s.ctx, tc.username)
		s.Assert().Equal(tc.expectErr, err)

		if tc.expectErr == nil {
			s.Assert().Equal(tc.username, u.Username)
			s.Assert().Equal(tc.password, u.Password)
		}
	}
}
