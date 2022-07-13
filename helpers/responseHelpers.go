package helpers

import (
	"fmt"
	"net/http"
)

func SendErrorResponse(responseWriter http.ResponseWriter, statusCode int, message string) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	responseWriter.Write(
		[]byte(
			fmt.Sprintf(`{
				"code": "%d"
				"message": "%s"
			}`, statusCode, message),
		),
	)
}

func SendIdentificationResponse(responseWriter http.ResponseWriter, isDog bool) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	if isDog {
		responseWriter.Write([]byte(`{
			"message": "Image contains a dog."
		}`))
	} else {
		responseWriter.Write([]byte(`{
			"message": "Image does not contain a dog."
		}`))
	}
}
