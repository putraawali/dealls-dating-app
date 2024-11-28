package usecases

import "github.com/sarulabs/di"

type Usecases struct {
	User  UserUsecase
	Swipe SwipeUsecase
}

func NewUsecase(di di.Container) *Usecases {
	return &Usecases{
		User:  NewUserUsecase(di),
		Swipe: NewSwipeUsecase(di),
	}
}
