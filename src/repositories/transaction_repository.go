package repositories

import (
	"context"
	"dealls-dating-app/src/constants"
	"dealls-dating-app/src/models"
	"dealls-dating-app/src/pkg/response"
	"errors"
	"net/http"

	"github.com/sarulabs/di"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Insert(ctx context.Context, transaction *models.Transaction) (err error)
	GetLatestTransactionByUserID(ctx context.Context, userID int64) (transaction models.Transaction, err error)
	Updates(ctx context.Context, transaction *models.Transaction) (err error)
}

type transactionRepository struct {
	db       *gorm.DB
	response *response.Response
}

func NewTransactionRepository(di di.Container) TransactionRepository {
	return &transactionRepository{
		db:       di.Get(constants.PG_DB).(*gorm.DB),
		response: di.Get(constants.RESPONSE).(*response.Response),
	}
}

func (t *transactionRepository) Insert(ctx context.Context, transaction *models.Transaction) (err error) {
	if err = t.db.Model(&models.Transaction{}).Create(transaction).WithContext(ctx).Error; err != nil {
		err = t.response.NewError().
			SetContext(ctx).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusInternalServerError)
	}
	return
}

func (t *transactionRepository) GetLatestTransactionByUserID(ctx context.Context, userID int64) (transaction models.Transaction, err error) {
	if err = t.db.Model(&models.Transaction{}).Last(&transaction, "user_id = ? AND status = ?", userID, "pending").WithContext(ctx).Error; err != nil {
		code := http.StatusInternalServerError
		msg := err

		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
			msg = errors.New("no transaction found")
		}

		err = t.response.NewError().
			SetContext(ctx).
			SetDetail(msg.Error()).
			SetMessage(msg).
			SetStatusCode(code)
	}
	return
}

func (t *transactionRepository) Updates(ctx context.Context, transaction *models.Transaction) (err error) {
	if err = t.db.Model(transaction).Updates(transaction).WithContext(ctx).Error; err != nil {
		err = t.response.NewError().
			SetContext(ctx).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusInternalServerError)
	}

	return
}
