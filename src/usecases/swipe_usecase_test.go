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
	swipeUsecaseTest struct {
		suite.Suite
		usecase usecases.SwipeUsecase
	}

	swipeUsecaseData struct {
		mockRepo *repositories.Repositories
		builder  *di.Builder
	}
)

func TestSwipeUsecase(t *testing.T) {
	suite.Run(t, new(swipeUsecaseTest))
}

func setupSwipeUsecase(t *testing.T) swipeUsecaseData {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockRepository(ctrl)

	mockDI := src_mock.Dependencies{
		Repository: mockRepo,
	}

	return swipeUsecaseData{
		builder:  src_mock.NewMockDependencies(mockDI),
		mockRepo: mockRepo,
	}
}

// Technical Debt of unit test
// Only create unit test for one function per usecase

func (s *swipeUsecaseTest) TestSwipePartner() {
	ctx := context.WithValue(context.TODO(), "request-id", "213")
	setup := setupSwipeUsecase(s.T())

	s.usecase = usecases.NewSwipeUsecase(setup.builder.Build())
	userRepo := setup.mockRepo.User.(*mock_repositories.MockUserRepository)
	swipeRepo := setup.mockRepo.Swipe.(*mock_repositories.MockSwipeRepository)

	paramLike := dtos.SwipePartnerParams{
		UserID:       1,
		TargetUserID: 2,
		Action:       "like",
	}

	s.Run("Success: success swipe user non premium", func() {
		userRepo.EXPECT().FindByID(ctx, paramLike.UserID).Return(models.User{UserID: 1}, nil)
		swipeRepo.EXPECT().CountSwipeUserToday(ctx, paramLike.UserID).Return(int64(1), nil)
		swipeRepo.EXPECT().FindByUserIDAndTargetIDToday(ctx, paramLike.UserID, paramLike.TargetUserID).Return(models.Swipe{}, nil)

		swipeData := models.Swipe{
			UserID:       paramLike.UserID,
			TargetUserID: paramLike.TargetUserID,
			Like:         true,
		}
		swipeRepo.EXPECT().Insert(ctx, &swipeData).Return(nil)

		err := s.usecase.SwipePartner(ctx, paramLike)
		s.Nil(err, "error should be nil")
	})

	paramPass := dtos.SwipePartnerParams{
		UserID:       1,
		TargetUserID: 2,
		Action:       "pass",
	}

	s.Run("Success: success swipe user premium", func() {
		userRepo.EXPECT().FindByID(ctx, paramPass.UserID).Return(models.User{UserID: 1, IsPremium: true}, nil)
		swipeRepo.EXPECT().FindByUserIDAndTargetIDToday(ctx, paramPass.UserID, paramPass.TargetUserID).Return(models.Swipe{}, nil)

		swipeData := models.Swipe{
			UserID:       paramPass.UserID,
			TargetUserID: paramPass.TargetUserID,
			Pass:         true,
		}
		swipeRepo.EXPECT().Insert(ctx, &swipeData).Return(nil)

		err := s.usecase.SwipePartner(ctx, paramPass)
		s.Nil(err, "error should be nil")
	})

	s.Run("Failed: failed swipe target/partner user", func() {
		userRepo.EXPECT().FindByID(ctx, paramLike.UserID).Return(models.User{UserID: 1}, nil)
		swipeRepo.EXPECT().CountSwipeUserToday(ctx, paramLike.UserID).Return(int64(1), nil)
		swipeRepo.EXPECT().FindByUserIDAndTargetIDToday(ctx, paramLike.UserID, paramLike.TargetUserID).Return(models.Swipe{}, nil)

		swipeData := models.Swipe{
			UserID:       paramLike.UserID,
			TargetUserID: paramLike.TargetUserID,
			Like:         true,
		}
		swipeRepo.EXPECT().Insert(ctx, &swipeData).Return(errors.New("mock error"))

		err := s.usecase.SwipePartner(ctx, paramLike)
		s.NotNil(err, "error should be not nil")
	})

	s.Run("Failed: failed swipe same partner user", func() {
		userRepo.EXPECT().FindByID(ctx, paramLike.UserID).Return(models.User{UserID: 1}, nil)
		swipeRepo.EXPECT().CountSwipeUserToday(ctx, paramLike.UserID).Return(int64(1), nil)
		swipeRepo.EXPECT().FindByUserIDAndTargetIDToday(ctx, paramLike.UserID, paramLike.TargetUserID).Return(models.Swipe{SwipeID: 1}, nil)

		err := s.usecase.SwipePartner(ctx, paramLike)
		s.NotNil(err, "error should be not nil")
	})

	s.Run("Failed: already swipe 10 times", func() {
		userRepo.EXPECT().FindByID(ctx, paramLike.UserID).Return(models.User{UserID: 1}, nil)
		swipeRepo.EXPECT().CountSwipeUserToday(ctx, paramLike.UserID).Return(int64(10), nil)

		err := s.usecase.SwipePartner(ctx, paramLike)
		s.NotNil(err, "error should be not nil")
	})

	s.Run("Failed: error get max attempt of swipe partner", func() {
		userRepo.EXPECT().FindByID(ctx, paramLike.UserID).Return(models.User{UserID: 1}, nil)
		swipeRepo.EXPECT().CountSwipeUserToday(ctx, paramLike.UserID).Return(int64(0), errors.New("mock error"))

		err := s.usecase.SwipePartner(ctx, paramLike)
		s.NotNil(err, "error should be not nil")
	})

	s.Run("Failed: error get max attempt of swipe partner", func() {
		userRepo.EXPECT().FindByID(ctx, paramLike.UserID).Return(models.User{}, errors.New("mock error"))

		err := s.usecase.SwipePartner(ctx, paramLike)
		s.NotNil(err, "error should be not nil")
	})
}
