package src_mock

import (
	"dealls-dating-app/src/constants"
	"dealls-dating-app/src/pkg/response"
	"dealls-dating-app/src/repositories"

	"github.com/sarulabs/di"
	"gorm.io/gorm"
)

type Dependencies struct {
	Repository *repositories.Repositories
	Postgres   *gorm.DB
	Mysql      *gorm.DB
}

func NewMockDependencies(d Dependencies) di.Container {
	builder, _ := di.NewBuilder()

	builder.Add(
		di.Def{
			Name: constants.RESPONSE,
			Build: func(ctn di.Container) (interface{}, error) {
				return response.NewResponse(), nil
			},
		},
		di.Def{
			Name: constants.REPOSITORY,
			Build: func(ctn di.Container) (interface{}, error) {
				return d.Repository, nil
			},
		},
		di.Def{
			Name: constants.PG_DB,
			Build: func(ctn di.Container) (interface{}, error) {
				return d.Postgres, nil
			},
		},
	)

	return builder.Build()
}
