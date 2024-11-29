package src

import (
	"dealls-dating-app/src/controllers"

	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
)

func NewRoutes(app *echo.Echo, di di.Container) {
	ctrl := controllers.NewController(di)

	UserRoutes(app, ctrl)
	SwipeRoutes(app, ctrl)
	TransactionRoutes(app, ctrl)
}

func UserRoutes(app *echo.Echo, ctrl *controllers.Controllers) {
	app.POST("/register", ctrl.User.Register)
	app.POST("/login", ctrl.User.Login)
	app.POST("/verify-email", ctrl.User.VerifyEmail)
}

func SwipeRoutes(app *echo.Echo, ctrl *controllers.Controllers) {
	app.GET("/available-partner", ctrl.Swipe.GetAvailablePartner)
	app.POST("/swipe-partner", ctrl.Swipe.SwipePartner)
}

func TransactionRoutes(app *echo.Echo, ctrl *controllers.Controllers) {
	app.POST("/request-payment-premium", ctrl.Transaction.InitTransaction)
	app.POST("/accept-transaction", ctrl.Transaction.AcceptTransaction)
}
