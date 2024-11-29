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
