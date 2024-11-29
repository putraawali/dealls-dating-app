package usecases

import (
	"context"
	"dealls-dating-app/src/constants"
	"dealls-dating-app/src/dtos"
	"dealls-dating-app/src/models"
	"dealls-dating-app/src/pkg/response"
	"dealls-dating-app/src/repositories"
	"math/rand"

	"github.com/sarulabs/di"
)

type TransactionUsecase interface {
	InitTransaction(ctx context.Context, param dtos.InitTransactionParam) (data dtos.InitTransactionResponse, err error)
	AcceptTransaction(ctx context.Context, param dtos.AcceptTransactionParam) (err error)
}

type transactionUsecase struct {
	repo     *repositories.Repositories
	response *response.Response
}

func NewTransactionUsecase(di di.Container) TransactionUsecase {
	return &transactionUsecase{
		repo:     di.Get(constants.REPOSITORY).(*repositories.Repositories),
		response: di.Get(constants.RESPONSE).(*response.Response),
	}
}

func (t *transactionUsecase) InitTransaction(ctx context.Context, param dtos.InitTransactionParam) (data dtos.InitTransactionResponse, err error) {
	// Mock process get VA Number for payment
	b := make([]rune, 15)

	letterRunes := []rune("0123456789")
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	data.VANumber = string(b)

	err = t.repo.Transaction.Insert(ctx, &models.Transaction{
		UserID:        param.UserID,
		PaymentMethod: param.PaymentMethod,
		VaNumber:      data.VANumber,
	})
	if err != nil {
		return
	}

	return
}

func (t *transactionUsecase) AcceptTransaction(ctx context.Context, param dtos.AcceptTransactionParam) (err error) {
	trx, err := t.repo.Transaction.GetLatestTransactionByUserID(ctx, param.UserID)
	if err != nil {
		return
	}

	trx.Status = "paid"

	err = t.repo.Transaction.Updates(ctx, &trx)
	if err != nil {
		return
	}

	err = t.repo.User.ActivatePremium(ctx, param.UserID)
	if err != nil {
		return
	}

	return
}
