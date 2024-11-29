package usecases_test

import (
	"context"
	"dealls-dating-app/src/dtos"
	src_mock "dealls-dating-app/src/mocks"
	"dealls-dating-app/src/models"
	"dealls-dating-app/src/repositories"
	mock_repositories "dealls-dating-app/src/repositories/mocks"
	"dealls-dating-app/src/usecases"
	"errors"
	"testing"

	"github.com/sarulabs/di"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type (
	userUsecaseTest struct {
		suite.Suite
		usecase usecases.UserUsecase
	}

	userUsecaseData struct {
		mockRepo *repositories.Repositories
		builder  *di.Builder
	}
)

func TestUserUsecase(t *testing.T) {
	suite.Run(t, new(userUsecaseTest))
}

func setupUserUsecase(t *testing.T) userUsecaseData {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockRepository(ctrl)

	mockDI := src_mock.Dependencies{
		Repository: mockRepo,
	}

	return userUsecaseData{
		builder:  src_mock.NewMockDependencies(mockDI),
		mockRepo: mockRepo,
	}
}

// Technical Debt of unit test
// Only create unit test for one function per usecase

func (u *userUsecaseTest) TestRegister() {
	ctx := context.WithValue(context.TODO(), "request-id", "213")
	setup := setupUserUsecase(u.T())

	u.usecase = usecases.NewUserUsecase(setup.builder.Build())

	userRepo := setup.mockRepo.User.(*mock_repositories.MockUserRepository)

	param := dtos.RegisterParam{
		Email:     "mail@mail.com",
		Password:  "Test1234",
		Sex:       "male",
		FirstName: "Dummy",
		LastName:  "User",
	}

	u.Run("Success: success register user", func() {
		userRepo.EXPECT().FindByEmail(ctx, param.Email).Return(models.User{}, nil)

		data := models.User{}
		data.RegisterToModel(param)
		userRepo.EXPECT().Insert(ctx, &data).Return(nil)

		err := u.usecase.Register(ctx, param)
		u.Nil(err, "Error should be nil")
	})

	u.Run("Failed: error insert new user", func() {
		userRepo.EXPECT().FindByEmail(ctx, param.Email).Return(models.User{}, nil)

		data := models.User{}
		data.RegisterToModel(param)
		userRepo.EXPECT().Insert(ctx, &data).Return(errors.New("mock error"))

		err := u.usecase.Register(ctx, param)
		u.NotNil(err, "Error should be not nil")
	})

	u.Run("Failed: email already used", func() {
		userRepo.EXPECT().FindByEmail(ctx, param.Email).Return(models.User{
			UserID: 1,
		}, nil)

		err := u.usecase.Register(ctx, param)
		u.NotNil(err, "Error should be not nil")
	})

	u.Run("Failed: error get user by email", func() {
		userRepo.EXPECT().FindByEmail(ctx, param.Email).Return(models.User{}, errors.New("mock error"))

		err := u.usecase.Register(ctx, param)
		u.NotNil(err, "Error should be not nil")
	})
}
