package controllers

import (
	"dealls-dating-app/src/constants"
	"dealls-dating-app/src/dtos"
	"dealls-dating-app/src/pkg/jwt"
	"dealls-dating-app/src/pkg/response"
	"dealls-dating-app/src/usecases"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
)

type SwipeController interface {
	GetAvailablePartner(echo.Context) error
	SwipePartner(echo.Context) error
}

type swipeController struct {
	uc       *usecases.Usecases
	response *response.Response
}

func NewSwipeController(di di.Container) SwipeController {
	return &swipeController{
		uc:       di.Get(constants.USECASE).(*usecases.Usecases),
		response: di.Get(constants.RESPONSE).(*response.Response),
	}
}

func (s *swipeController) GetAvailablePartner(c echo.Context) error {
	userData := jwt.GetUserData(c)

	data, err := s.uc.Swipe.GetAvailablePartner(c.Request().Context(), userData.UserID)
	if err != nil {
		return c.JSON(s.response.Send(0, nil, err))
	}

	return c.JSON(s.response.Send(http.StatusOK, data, nil))
}

func (s *swipeController) SwipePartner(c echo.Context) error {
	userData := jwt.GetUserData(c)

	param := dtos.SwipePartnerParams{}

	err := c.Bind(&param)
	if err != nil {
		err = s.response.NewError().
			SetContext(c.Request().Context()).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusBadRequest)

		return c.JSON(s.response.Send(0, nil, err))
	}

	param.UserID = userData.UserID

	err = param.Validate()
	if err != nil {
		err = s.response.NewError().
			SetContext(c.Request().Context()).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusBadRequest)

		return c.JSON(s.response.Send(0, nil, err))
	}

	err = s.uc.Swipe.SwipePartner(c.Request().Context(), param)
	if err != nil {
		return c.JSON(s.response.Send(0, nil, err))
	}

	return c.JSON(s.response.Send(http.StatusOK, nil, nil))
}
