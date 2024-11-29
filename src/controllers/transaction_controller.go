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

type TransactionController interface {
	InitTransaction(echo.Context) error
	AcceptTransaction(echo.Context) error
}

type transactionController struct {
	uc       *usecases.Usecases
	response *response.Response
}

func NewTransactionController(di di.Container) TransactionController {
	return &transactionController{
		uc:       di.Get(constants.USECASE).(*usecases.Usecases),
		response: di.Get(constants.RESPONSE).(*response.Response),
	}
}

func (t *transactionController) InitTransaction(c echo.Context) error {
	userData := jwt.GetUserData(c)

	param := dtos.InitTransactionParam{}

	err := c.Bind(&param)
	if err != nil {
		err = t.response.NewError().
			SetContext(c.Request().Context()).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusBadRequest)

		return c.JSON(t.response.Send(0, nil, err))
	}

	param.UserID = userData.UserID

	data, err := t.uc.Transaction.InitTransaction(c.Request().Context(), param)
	if err != nil {
		return c.JSON(t.response.Send(0, nil, err))
	}

	return c.JSON(t.response.Send(http.StatusCreated, data, nil))
}

func (t *transactionController) AcceptTransaction(c echo.Context) error {
	userData := jwt.GetUserData(c)

	param := dtos.AcceptTransactionParam{}

	err := c.Bind(&param)
	if err != nil {
		err = t.response.NewError().
			SetContext(c.Request().Context()).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusBadRequest)

		return c.JSON(t.response.Send(0, nil, err))
	}

	param.UserID = userData.UserID

	err = t.uc.Transaction.AcceptTransaction(c.Request().Context(), param)
	if err != nil {
		return c.JSON(t.response.Send(0, nil, err))
	}

	return c.JSON(t.response.Send(http.StatusOK, nil, nil))
}
