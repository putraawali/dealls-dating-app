package repositories_test

import (
	"context"
	src_mock "dealls-dating-app/src/mocks"
	"dealls-dating-app/src/models"
	mock_connections "dealls-dating-app/src/pkg/connections/mocks"
	"dealls-dating-app/src/repositories"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type userRepositoryTest struct {
	suite.Suite
	repoPostgres repositories.UserRepository
	postgres     *gorm.DB
	mockPostgres sqlmock.Sqlmock
	ctx          context.Context
}

var columnUser = []string{
	"user_id",
	"email",
	"is_verified",
	"sex",
	"first_name",
	"last_name",
	"password",
	"created_at",
	"updated_at",
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(userRepositoryTest))
}

// Before each test
func (u *userRepositoryTest) SetupTest() {
	u.postgres, u.mockPostgres = mock_connections.NewMockPostgresConnection()

	diPostgres := src_mock.NewMockDependencies(src_mock.Dependencies{
		Postgres: u.postgres,
	})

	u.repoPostgres = repositories.NewUserRepository(diPostgres.Build())

	u.ctx = context.WithValue(context.TODO(), "request-id", "213")
}

func (u *userRepositoryTest) TestInsert() {
	u.Run("Success insert user", func() {
		data := models.User{
			Email:     "email@mail.com",
			Sex:       "male",
			FirstName: "First Name",
			LastName:  "Last Name",
			Password:  "pw-123",
		}

		query := "INSERT INTO \"users\" "

		u.mockPostgres.ExpectBegin()

		u.mockPostgres.ExpectQuery(query).
			WithArgs(
				data.Email,
				data.IsVerified,
				data.Sex,
				data.FirstName,
				data.LastName,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

		u.mockPostgres.ExpectCommit()

		err := u.repoPostgres.Insert(u.ctx, &data)
		u.Nil(err, "error should be nil")
	})

	u.Run("Failed insert user", func() {
		data := models.User{
			Email:     "email@mail.com",
			FirstName: "First Name",
			LastName:  "Last Name",
			Sex:       "female",
			Password:  "pw-123",
		}

		query := "INSERT INTO \"users\" "

		u.mockPostgres.ExpectBegin()

		u.mockPostgres.ExpectQuery(query).
			WithArgs(
				data.Email,
				data.IsVerified,
				data.Sex,
				data.FirstName,
				data.LastName,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnError(errors.New("mock error"))

		u.mockPostgres.ExpectRollback()

		err := u.repoPostgres.Insert(u.ctx, &data)
		u.NotNil(err, "error should be not nil")
	})

	u.Run("Failed insert user wrong type of sex", func() {
		data := models.User{
			Email:     "email@mail.com",
			FirstName: "First Name",
			LastName:  "Last Name",
			Sex:       "",
			Password:  "pw-123",
		}

		u.mockPostgres.ExpectBegin()

		u.mockPostgres.ExpectRollback()

		err := u.repoPostgres.Insert(u.ctx, &data)
		u.NotNil(err, "error should be not nil")
	})
}

func (u *userRepositoryTest) TestFindByEmail() {
	u.Run("Success find user", func() {
		data := models.User{
			Email:     "email@mail.com",
			Sex:       "male",
			FirstName: "First Name",
			LastName:  "Last Name",
			Password:  "pw-123",
		}

		query := "SELECT * FROM \"users\" "

		u.mockPostgres.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(data.Email, 1).
			WillReturnRows(
				sqlmock.NewRows(columnUser).
					AddRow(
						data.UserID,
						data.Email,
						data.IsVerified,
						data.Sex,
						data.FirstName,
						data.LastName,
						data.Password,
						data.CreatedAt,
						data.UpdatedAt,
					),
			)

		user, err := u.repoPostgres.FindByEmail(u.ctx, data.Email)
		u.Nil(err, "error should be nil")
		u.Equal(user, data, "User with data should be equal")
	})

	u.Run("Failed user not found", func() {
		query := "SELECT * FROM \"users\" "

		u.mockPostgres.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(gorm.ErrRecordNotFound)

		user, err := u.repoPostgres.FindByEmail(u.ctx, "email@mail.com")
		u.NotNil(err, "error should be nil")
		u.Empty(user, "User should be empty")
	})

	u.Run("Failed error internal server error", func() {
		query := "SELECT * FROM \"users\" "

		u.mockPostgres.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("mock error"))

		user, err := u.repoPostgres.FindByEmail(u.ctx, "email@mail.com")
		u.NotNil(err, "error should be nil")
		u.Empty(user, "User should be empty")
	})
}

func (u *userRepositoryTest) TestVerifyEmail() {
	u.Run("Success verify email user", func() {
		query := "UPDATE \"users\" SET "

		u.mockPostgres.ExpectBegin()

		u.mockPostgres.ExpectExec(query).
			WithArgs(
				true,
				sqlmock.AnyArg(),
				"email@mail.com",
			).WillReturnResult(sqlmock.NewResult(1, 1))

		u.mockPostgres.ExpectCommit()

		err := u.repoPostgres.VerifyEmail(u.ctx, "email@mail.com")
		u.Nil(err, "error should be nil")
	})

	u.Run("Failed verify email user", func() {
		query := "UPDATE \"users\" SET "

		u.mockPostgres.ExpectBegin()

		u.mockPostgres.ExpectExec(query).
			WithArgs(
				true,
				sqlmock.AnyArg(),
				"email@mail.com",
			).WillReturnError(errors.New("mock error"))

		u.mockPostgres.ExpectCommit()

		err := u.repoPostgres.VerifyEmail(u.ctx, "email@mail.com")
		u.NotNil(err, "error should be not nil")
	})
}
