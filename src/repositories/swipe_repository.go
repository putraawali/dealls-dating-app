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

type SwipeRepository interface {
	FindByUserIDAndTargetIDToday(ctx context.Context, userID, targetID int64) (data models.Swipe, err error)
	Insert(ctx context.Context, swipe *models.Swipe) (err error)
	CountSwipeUserToday(ctx context.Context, userID int64) (count int64, err error)
}

type swipeRepository struct {
	db       *gorm.DB
	response *response.Response
}

func NewSwipeRepository(di di.Container) SwipeRepository {
	return &swipeRepository{
		db:       di.Get(constants.PG_DB).(*gorm.DB),
		response: di.Get(constants.RESPONSE).(*response.Response),
	}
}

func (s *swipeRepository) FindByUserIDAndTargetIDToday(ctx context.Context, userID, targetID int64) (data models.Swipe, err error) {
	if err = s.db.Model(&models.Swipe{}).
		First(
			&data,
			"user_id = ? AND target_user_id = ? AND DATE(created_at) = DATE(now())", userID, targetID,
		).WithContext(ctx).Error; err != nil {
		code := http.StatusInternalServerError
		msg := err

		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
			msg = errors.New("user not found")
		}

		err = s.response.NewError().
			SetContext(ctx).
			SetDetail(msg.Error()).
			SetMessage(msg).
			SetStatusCode(code)
	}

	return
}

func (s *swipeRepository) Insert(ctx context.Context, swipe *models.Swipe) (err error) {
	if err = s.db.Model(&models.Swipe{}).Create(swipe).WithContext(ctx).Error; err != nil {
		err = s.response.NewError().
			SetContext(ctx).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusInternalServerError)
	}

	return
}

func (s *swipeRepository) CountSwipeUserToday(ctx context.Context, userID int64) (count int64, err error) {
	if err = s.db.Model(&models.Swipe{}).Where("user_id = ? AND DATE(created_at) = DATE(now())", userID).Count(&count).WithContext(ctx).Error; err != nil {
		err = s.response.NewError().
			SetContext(ctx).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusInternalServerError)
	}

	return
}
