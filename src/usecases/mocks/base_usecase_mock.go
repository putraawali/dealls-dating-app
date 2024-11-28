package mock_usecases

import (
	"dealls-dating-app/src/usecases"

	"go.uber.org/mock/gomock"
)

func NewBaseUsecaseMock(ctrl *gomock.Controller) *usecases.Usecases {
	return &usecases.Usecases{
		User:  NewMockUserUsecase(ctrl),
		Swipe: NewMockSwipeUsecase(ctrl),
	}
}
