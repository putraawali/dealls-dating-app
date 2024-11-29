package controllers

import "github.com/sarulabs/di"

type Controllers struct {
	User        UserController
	Swipe       SwipeController
	Transaction TransactionController
}

func NewController(di di.Container) *Controllers {
	return &Controllers{
		User:        NewUserController(di),
		Swipe:       NewSwipeController(di),
		Transaction: NewTransactionController(di),
	}
}
