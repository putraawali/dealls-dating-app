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
	transactionUsecaseTest struct {
		suite.Suite
		usecase usecases.TransactionUsecase
	}

	transactionUsecaseData struct {
		mockRepo *repositories.Repositories
		builder  *di.Builder
	}
)

func TestTranactionUsecase(t *testing.T) {
	suite.Run(t, new(transactionUsecaseTest))
}

func setupTransactionUsecase(t *testing.T) transactionUsecaseData {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockRepository(ctrl)

	mockDI := src_mock.Dependencies{
		Repository: mockRepo,
	}

	return transactionUsecaseData{
		builder:  src_mock.NewMockDependencies(mockDI),
		mockRepo: mockRepo,
	}
}

// Technical Debt of unit test
// Only create unit test for one function per usecase

func (t *transactionUsecaseTest) TestAcceptTransaction() {
	ctx := context.WithValue(context.TODO(), "request-id", "213")
	setup := setupTransactionUsecase(t.T())

	t.usecase = usecases.NewTransactionUsecase(setup.builder.Build())
	userRepo := setup.mockRepo.User.(*mock_repositories.MockUserRepository)
	transactionRepo := setup.mockRepo.Transaction.(*mock_repositories.MockTransactionRepository)

	param := dtos.AcceptTransactionParam{
		PaymentMethod: "virtual_account",
		VaNumber:      "1234567890",
		UserID:        1,
	}

	t.Run("Success: success accept transaction", func() {
		latestTrx := models.Transaction{
			TransactionID: 1,
			Status:        "pending",
			PaymentMethod: "virtual_account",
			VaNumber:      "1234567890",
		}
		transactionRepo.EXPECT().GetLatestTransactionByUserID(ctx, param.UserID).Return(latestTrx, nil)

		latestTrx.Status = "paid"
		transactionRepo.EXPECT().Updates(ctx, &latestTrx).Return(nil)

		userRepo.EXPECT().ActivatePremium(ctx, param.UserID).Return(nil)

		err := t.usecase.AcceptTransaction(ctx, param)
		t.Nil(err, "error should be nil")
	})

	t.Run("Failed: failed update user premium status", func() {
		latestTrx := models.Transaction{
			TransactionID: 1,
			Status:        "pending",
			PaymentMethod: "virtual_account",
			VaNumber:      "1234567890",
		}
		transactionRepo.EXPECT().GetLatestTransactionByUserID(ctx, param.UserID).Return(latestTrx, nil)

		latestTrx.Status = "paid"
		transactionRepo.EXPECT().Updates(ctx, &latestTrx).Return(nil)

		userRepo.EXPECT().ActivatePremium(ctx, param.UserID).Return(errors.New("mock error"))

		err := t.usecase.AcceptTransaction(ctx, param)
		t.NotNil(err, "error should be not nil")
	})

	t.Run("Failed: failed update transaction status to paid", func() {
		latestTrx := models.Transaction{
			TransactionID: 1,
			Status:        "pending",
			PaymentMethod: "virtual_account",
			VaNumber:      "1234567890",
		}
		transactionRepo.EXPECT().GetLatestTransactionByUserID(ctx, param.UserID).Return(latestTrx, nil)

		latestTrx.Status = "paid"
		transactionRepo.EXPECT().Updates(ctx, &latestTrx).Return(errors.New("mock error"))

		err := t.usecase.AcceptTransaction(ctx, param)
		t.NotNil(err, "error should be not nil")
	})

	t.Run("Failed: failed get latest transaction", func() {
		transactionRepo.EXPECT().GetLatestTransactionByUserID(ctx, param.UserID).Return(models.Transaction{}, errors.New("mock error"))

		err := t.usecase.AcceptTransaction(ctx, param)
		t.NotNil(err, "error should be not nil")
	})
}

// func (s *transactionUsecaseTest) TestSwipePartner() {
// 	swipeRepo := setup.mockRepo.Swipe.(*mock_repositories.MockSwipeRepository)

// 	paramLike := dtos.SwipePartnerParams{
// 		UserID:       1,
// 		TargetUserID: 2,
// 		Action:       "like",
// 	}

// 	s.Run("Success: success swipe user non premium", func() {
// 		userRepo.EXPECT().FindByID(ctx, paramLike.UserID).Return(models.User{UserID: 1}, nil)
// 		swipeRepo.EXPECT().CountSwipeUserToday(ctx, paramLike.UserID).Return(int64(1), nil)
// 		swipeRepo.EXPECT().FindByUserIDAndTargetIDToday(ctx, paramLike.UserID, paramLike.TargetUserID).Return(models.Swipe{}, nil)

// 		swipeData := models.Swipe{
// 			UserID:       paramLike.UserID,
// 			TargetUserID: paramLike.TargetUserID,
// 			Like:         true,
// 		}
// 		swipeRepo.EXPECT().Insert(ctx, &swipeData).Return(nil)

// 		err := s.usecase.SwipePartner(ctx, paramLike)
// 		s.Nil(err, "error should be nil")
// 	})

// 	paramPass := dtos.SwipePartnerParams{
// 		UserID:       1,
// 		TargetUserID: 2,
// 		Action:       "pass",
// 	}

// 	s.Run("Success: success swipe user premium", func() {
// 		userRepo.EXPECT().FindByID(ctx, paramPass.UserID).Return(models.User{UserID: 1, IsPremium: true}, nil)
// 		swipeRepo.EXPECT().FindByUserIDAndTargetIDToday(ctx, paramPass.UserID, paramPass.TargetUserID).Return(models.Swipe{}, nil)

// 		swipeData := models.Swipe{
// 			UserID:       paramPass.UserID,
// 			TargetUserID: paramPass.TargetUserID,
// 			Pass:         true,
// 		}
// 		swipeRepo.EXPECT().Insert(ctx, &swipeData).Return(nil)

// 		err := s.usecase.SwipePartner(ctx, paramPass)
// 		s.Nil(err, "error should be nil")
// 	})

// 	s.Run("Failed: failed swipe target/partner user", func() {
// 		userRepo.EXPECT().FindByID(ctx, paramLike.UserID).Return(models.User{UserID: 1}, nil)
// 		swipeRepo.EXPECT().CountSwipeUserToday(ctx, paramLike.UserID).Return(int64(1), nil)
// 		swipeRepo.EXPECT().FindByUserIDAndTargetIDToday(ctx, paramLike.UserID, paramLike.TargetUserID).Return(models.Swipe{}, nil)

// 		swipeData := models.Swipe{
// 			UserID:       paramLike.UserID,
// 			TargetUserID: paramLike.TargetUserID,
// 			Like:         true,
// 		}
// 		swipeRepo.EXPECT().Insert(ctx, &swipeData).Return(errors.New("mock error"))

// 		err := s.usecase.SwipePartner(ctx, paramLike)
// 		s.NotNil(err, "error should be not nil")
// 	})

// 	s.Run("Failed: failed swipe same partner user", func() {
// 		userRepo.EXPECT().FindByID(ctx, paramLike.UserID).Return(models.User{UserID: 1}, nil)
// 		swipeRepo.EXPECT().CountSwipeUserToday(ctx, paramLike.UserID).Return(int64(1), nil)
// 		swipeRepo.EXPECT().FindByUserIDAndTargetIDToday(ctx, paramLike.UserID, paramLike.TargetUserID).Return(models.Swipe{SwipeID: 1}, nil)

// 		err := s.usecase.SwipePartner(ctx, paramLike)
// 		s.NotNil(err, "error should be not nil")
// 	})

// 	s.Run("Failed: already swipe 10 times", func() {
// 		userRepo.EXPECT().FindByID(ctx, paramLike.UserID).Return(models.User{UserID: 1}, nil)
// 		swipeRepo.EXPECT().CountSwipeUserToday(ctx, paramLike.UserID).Return(int64(10), nil)

// 		err := s.usecase.SwipePartner(ctx, paramLike)
// 		s.NotNil(err, "error should be not nil")
// 	})

// 	s.Run("Failed: error get max attempt of swipe partner", func() {
// 		userRepo.EXPECT().FindByID(ctx, paramLike.UserID).Return(models.User{UserID: 1}, nil)
// 		swipeRepo.EXPECT().CountSwipeUserToday(ctx, paramLike.UserID).Return(int64(0), errors.New("mock error"))

// 		err := s.usecase.SwipePartner(ctx, paramLike)
// 		s.NotNil(err, "error should be not nil")
// 	})

// 	s.Run("Failed: error get max attempt of swipe partner", func() {
// 		userRepo.EXPECT().FindByID(ctx, paramLike.UserID).Return(models.User{}, errors.New("mock error"))

// 		err := s.usecase.SwipePartner(ctx, paramLike)
// 		s.NotNil(err, "error should be not nil")
// 	})
// }
