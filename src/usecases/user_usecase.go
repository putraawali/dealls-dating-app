package usecases

import (
	"context"
	"dealls-dating-app/src/constants"
	"dealls-dating-app/src/dtos"
	"dealls-dating-app/src/models"
	"dealls-dating-app/src/pkg/helpers"
	"dealls-dating-app/src/pkg/jwt"
	"dealls-dating-app/src/pkg/response"
	"dealls-dating-app/src/repositories"
	"errors"
	"fmt"
	"net/http"

	"github.com/sarulabs/di"
)

type UserUsecase interface {
	Register(ctx context.Context, data dtos.RegisterParam) (err error)
	Login(ctx context.Context, data dtos.LoginParam) (response dtos.LoginResponse, err error)
	VerifyEmail(ctx context.Context, data dtos.VerifyEmailParam) (err error)
}

type userUsecase struct {
	repo     *repositories.Repositories
	response *response.Response
}

func NewUserUsecase(di di.Container) UserUsecase {
	return &userUsecase{
		repo:     di.Get(constants.REPOSITORY).(*repositories.Repositories),
		response: di.Get(constants.RESPONSE).(*response.Response),
	}
}

func (u *userUsecase) Register(ctx context.Context, data dtos.RegisterParam) (err error) {
	user := models.User{}
	user.RegisterToModel(data)

	userData, err := u.repo.User.FindByEmail(ctx, data.Email)
	if err != nil {
		if !helpers.IsErrorNotFound(err) {
			return
		}
	}

	if userData.UserID != 0 {
		err = response.NewResponse().NewError().
			SetContext(ctx).
			SetDetail(fmt.Sprintf("email %s already used", data.Email)).
			SetMessage(fmt.Errorf("email %s already used", data.Email)).
			SetStatusCode(http.StatusBadRequest)

		return
	}

	return u.repo.User.Insert(ctx, &user)
}

func (u *userUsecase) Login(ctx context.Context, data dtos.LoginParam) (response dtos.LoginResponse, err error) {
	user, err := u.repo.User.FindByEmail(ctx, data.Email)
	if err != nil {
		return
	}

	if !helpers.ComparePassword([]byte(user.Password), []byte(data.Password)) {
		err = u.response.NewError().
			SetContext(ctx).
			SetMessage(errors.New("wrong password")).
			SetDetail("wrong password").
			SetStatusCode(http.StatusUnauthorized)
		return
	}

	response.AccessToken = jwt.GenerateToken(user.UserID, user.Email)

	return
}

func (u *userUsecase) VerifyEmail(ctx context.Context, data dtos.VerifyEmailParam) (err error) {
	user, err := u.repo.User.FindByEmail(ctx, data.Email)
	if err != nil {
		return
	}

	return u.repo.User.VerifyEmail(ctx, user.Email)
}
