package controllers

import (
	"dealls-dating-app/src/constants"
	"dealls-dating-app/src/dtos"
	"dealls-dating-app/src/pkg/response"
	"dealls-dating-app/src/usecases"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
)

type UserController interface {
	Register(echo.Context) error
	Login(echo.Context) error
	VerifyEmail(echo.Context) error
}

type userController struct {
	uc       *usecases.Usecases
	response *response.Response
}

func NewUserController(di di.Container) UserController {
	return &userController{
		uc:       di.Get(constants.USECASE).(*usecases.Usecases),
		response: di.Get(constants.RESPONSE).(*response.Response),
	}
}

func (u *userController) Register(c echo.Context) error {
	param := dtos.RegisterParam{}

	err := c.Bind(&param)
	if err != nil {
		err = u.response.NewError().
			SetContext(c.Request().Context()).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusBadRequest)

		return c.JSON(u.response.Send(0, nil, err))
	}

	err = param.Validate()
	if err != nil {
		err = u.response.NewError().
			SetContext(c.Request().Context()).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusBadRequest)

		return c.JSON(u.response.Send(0, nil, err))
	}

	err = u.uc.User.Register(c.Request().Context(), param)
	if err != nil {
		return c.JSON(u.response.Send(0, nil, err))
	}

	return c.JSON(u.response.Send(http.StatusCreated, nil, nil))
}

func (u *userController) Login(c echo.Context) error {
	param := dtos.LoginParam{}

	err := c.Bind(&param)
	if err != nil {
		err = u.response.NewError().
			SetContext(c.Request().Context()).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusBadRequest)

		return c.JSON(u.response.Send(0, nil, err))
	}

	err = param.Validate()
	if err != nil {
		err = u.response.NewError().
			SetContext(c.Request().Context()).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusBadRequest)

		return c.JSON(u.response.Send(0, nil, err))
	}

	data, err := u.uc.User.Login(c.Request().Context(), param)
	if err != nil {
		return c.JSON(u.response.Send(0, nil, err))
	}

	return c.JSON(u.response.Send(http.StatusOK, data, nil))
}

func (u *userController) VerifyEmail(c echo.Context) error {
	param := dtos.VerifyEmailParam{}

	err := c.Bind(&param)
	if err != nil {
		err = u.response.NewError().
			SetContext(c.Request().Context()).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusBadRequest)

		return c.JSON(u.response.Send(0, nil, err))
	}

	err = param.Validate()
	if err != nil {
		err = u.response.NewError().
			SetContext(c.Request().Context()).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusBadRequest)

		return c.JSON(u.response.Send(0, nil, err))
	}

	err = u.uc.User.VerifyEmail(c.Request().Context(), param)
	if err != nil {
		return c.JSON(u.response.Send(0, nil, err))
	}

	return c.JSON(u.response.Send(http.StatusOK, nil, nil))
}
