package usecases_test

import (
	"dealls-dating-app/src/repositories"
	"dealls-dating-app/src/usecases"

	"github.com/stretchr/testify/suite"
)

type (
	userUsecaseTest struct {
		suite.Suite
		usecase usecases.UserUsecase
	}

	userUsecaseData struct {
		mockRepo repositories.Repositories
	}
)
