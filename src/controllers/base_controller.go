package controllers

import "github.com/sarulabs/di"

type Controllers struct {
	User  UserController
	Swipe SwipeController
}

func NewController(di di.Container) *Controllers {
	return &Controllers{
		User:  NewUserController(di),
		Swipe: NewSwipeController(di),
	}
}
