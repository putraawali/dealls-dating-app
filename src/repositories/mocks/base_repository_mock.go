package mock_repositories

import (
	"dealls-dating-app/src/repositories"

	"go.uber.org/mock/gomock"
)

func NewMockRepository(ctrl *gomock.Controller) *repositories.Repositories {
	return &repositories.Repositories{
		User:        NewMockUserRepository(ctrl),
		Swipe:       NewMockSwipeRepository(ctrl),
		Transaction: NewMockTransactionRepository(ctrl),
	}
}
