package src

import (
	"dealls-dating-app/src/pkg/middlewares"

	"github.com/labstack/echo/v4"
)

type Module struct{}

func (m *Module) New(app *echo.Echo) {
	di := dependencyInjection()
	middlewares.UseMiddlwares(app, di)
	NewRoutes(app, di)
}
