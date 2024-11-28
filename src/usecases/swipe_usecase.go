package usecases

import (
	"context"
	"dealls-dating-app/src/constants"
	"dealls-dating-app/src/dtos"
	"dealls-dating-app/src/models"
	"dealls-dating-app/src/pkg/response"
	"dealls-dating-app/src/repositories"
	usecase_helpers "dealls-dating-app/src/usecases/helpers"
	"errors"
	"net/http"

	"github.com/sarulabs/di"
	"golang.org/x/sync/errgroup"
)

type SwipeUsecase interface {
	GetAvailablePartner(ctx context.Context, userId int64) (data []models.User, err error)
	SwipePartner(ctx context.Context, data dtos.SwipePartnerParams) (err error)
}

type swipeUsecase struct {
	repo     *repositories.Repositories
	response *response.Response
}

func NewSwipeUsecase(di di.Container) SwipeUsecase {
	return &swipeUsecase{
		repo:     di.Get(constants.REPOSITORY).(*repositories.Repositories),
		response: di.Get(constants.RESPONSE).(*response.Response),
	}
}

func (s *swipeUsecase) GetAvailablePartner(ctx context.Context, userId int64) (data []models.User, err error) {
	user, err := s.repo.User.FindByID(ctx, userId)
	if err != nil {
		return
	}

	interest := map[string]string{
		"male":   "female",
		"female": "male",
	}

	limit, offset := 10, 0

	tempData := usecase_helpers.NewHandleAvailableUser()

	var availableUser []models.User

	errs, _ := errgroup.WithContext(ctx)
	errs.SetLimit(10)

	// Look for available person to match
	for {
		availableUser, err = s.repo.User.FindBySex(ctx, interest[user.Sex], limit, offset)
		if err != nil {
			return
		}

		if len(availableUser) == 0 {
			break
		}

		// Check if user already pass / swipe today
		for i, userAvailable := range availableUser {
			errs.Go(func() error {
				swipeData, _ := s.repo.Swipe.FindByUserIDAndTargetIDToday(ctx, user.UserID, userAvailable.UserID)

				if swipeData.SwipeID == 0 {
					if tempData.GetCurrent() < 10 {
						tempData.Add(availableUser[i])
					}
				}

				return nil
			})
		}

		errs.Wait()

		if tempData.GetCurrent() == 10 {
			break
		} else {
			offset += 10
		}
	}

	data = tempData.FinalData()

	return
}

func (s *swipeUsecase) SwipePartner(ctx context.Context, data dtos.SwipePartnerParams) (err error) {
	count, err := s.repo.Swipe.CountSwipeUserToday(ctx, data.UserID)
	if err != nil {
		return
	}

	if count > 10 {
		msg := errors.New("anda telah mencapai batas maksimum swipe hari ini")
		err = s.response.NewError().
			SetContext(ctx).
			SetDetail(msg.Error()).
			SetMessage(msg).
			SetStatusCode(http.StatusBadRequest)
		return
	}

	currentSwipeData, _ := s.repo.Swipe.FindByUserIDAndTargetIDToday(ctx, data.UserID, data.TargetUserID)
	if currentSwipeData.SwipeID != 0 {
		msg := errors.New("anda telah melakukan swipe pada user tersebut hari ini")
		err = s.response.NewError().
			SetContext(ctx).
			SetDetail(msg.Error()).
			SetMessage(msg).
			SetStatusCode(http.StatusBadRequest)
		return
	}

	swipeData := models.Swipe{
		UserID:       data.UserID,
		TargetUserID: data.TargetUserID,
	}

	if data.Action == "pass" {
		swipeData.Pass = true
	} else if data.Action == "like" {
		swipeData.Like = true
	}

	return s.repo.Swipe.Insert(ctx, &swipeData)
}
