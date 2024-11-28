package src

import (
	"dealls-dating-app/src/constants"
	"dealls-dating-app/src/pkg/connections"
	"dealls-dating-app/src/pkg/response"
	"dealls-dating-app/src/repositories"
	"dealls-dating-app/src/usecases"

	"github.com/sarulabs/di"
)

func dependencyInjection() di.Container {
	builder, _ := di.NewBuilder()

	pg, err := connections.NewPostgreConnection()

	builder.Add(
		di.Def{
			Name: constants.RESPONSE,
			Build: func(ctn di.Container) (interface{}, error) {
				return response.NewResponse(), nil
			},
		},
		di.Def{
			Name: constants.PG_DB,
			Build: func(ctn di.Container) (interface{}, error) {
				return pg, err
			},
		},
		di.Def{
			Name: constants.REPOSITORY,
			Build: func(ctn di.Container) (interface{}, error) {
				return repositories.NewRepository(builder.Build()), nil
			},
		},
		di.Def{
			Name: constants.USECASE,
			Build: func(ctn di.Container) (interface{}, error) {
				return usecases.NewUsecase(builder.Build()), nil
			},
		},
	)

	return builder.Build()
}
