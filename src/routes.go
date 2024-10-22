package src

import (
	"go-boilerplate-v2/src/controllers"

	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
)

func NewRoutes(app *echo.Echo, di di.Container) {
	ctrl := controllers.NewController(di)

	app.POST("/register", ctrl.User.Register)
	app.POST("/login", ctrl.User.Login)

	// app.GET("/test", func(c echo.Context) error {
	// 	userData := jwt.GetUserData(c)
	// 	return c.JSON(200, userData)
	// })
}
