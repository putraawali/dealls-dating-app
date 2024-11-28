package helpers

import (
	"dealls-dating-app/src/pkg/response"
	"net/http"
)

func IsErrorNotFound(err error) bool {
	errData, ok := err.(*response.ErrorResponse)
	if ok {
		return errData.GetStatusCode() == http.StatusNotFound
	}

	return false
}
