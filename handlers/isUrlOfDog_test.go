package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIsUrlOfDogHandler(t *testing.T) {
	const successMessageWithDog = `"message": "Image contains a dog."`
	const successMessageWithoutDog = `"message": "Image does not contain a dog."`
	const invalidURLMessage = `"error": "URL is invalid."`
	const unusableURLMessage = `"error": "URL is invalid or not accessible."`
	const invalidImageMessage = `"error": "URL does not contain raw image data."`
	tests := []struct {
		name           string
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "test happy path with dog",
			in:             httptest.NewRequest("GET", "/isUrlOfDog?imageUrl=https://image.shutterstock.com/image-photo/golden-retriever-dog-lying-against-260nw-1093957301.jpg", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   successMessageWithDog,
		},
		{
			name:           "test happy path without dog",
			in:             httptest.NewRequest("GET", "/isUrlOfDog?imageUrl=https://cdn.pixabay.com/photo/2015/04/19/08/32/marguerite-729510_960_720.jpg", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   successMessageWithoutDog,
		},
		{
			name:           "test with invalid scheme",
			in:             httptest.NewRequest("GET", "/isUrlOfDog?imageUrl=about://www.google.com", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   invalidURLMessage,
		},
		{
			name:           "Test with valid URL not going anywhere",
			in:             httptest.NewRequest("GET", "/isUrlOfDog?imageUrl=https://www.googleidsuhfgoijk.com", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   unusableURLMessage,
		},
		{
			name:           "test with URL not containing just raw image data",
			in:             httptest.NewRequest("GET", "/isUrlOfDog?imageUrl=https://www.google.com", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   invalidImageMessage,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			IsUrlOfDog(test.out, test.in)
			if test.out.Code != test.expectedStatus {
				t.Logf("expected status code: %d\ngot: %d\n", test.expectedStatus, test.out.Code)
				t.Fail()
			}

			body := test.out.Body.String()
			if !strings.Contains(body, test.expectedBody) {
				t.Logf("expected in result: %s\ngot: %s\n", test.expectedBody, body)
				t.Fail()
			}
		})
	}
}
