package handlers

import (
	"bytes"
	"dogidentifier/helpers"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"github.com/hashicorp/go-retryablehttp"
)

func IsUrlOfDog(responseWriter http.ResponseWriter, request *http.Request) {
	googleApiKey := os.Getenv("IMAGE_IDENTIFIER_GOOGLE_API_KEY");

	if googleApiKey == "" {
		fmt.Println("Did not find valid Google API key!");
		responseMessage := "Could not handle request. Please try again later.";
		helpers.SendErrorResponse(responseWriter, http.StatusInternalServerError, responseMessage);
		return;
	}

	imageURL, err := verifyImageUrl(request);
	if err != nil {
		responseMessage := "URL format is invalid.";
		helpers.SendErrorResponse(responseWriter, http.StatusBadRequest, responseMessage);
		return;
	}

	response, err := getResultFromVisionApiWithURL(imageURL, googleApiKey);
	if err != nil {
		fmt.Printf("Failed to get response from API.");
		responseMessage := "Could not handle request. Please try again later.";
		helpers.SendErrorResponse(responseWriter, http.StatusBadRequest, responseMessage);
		return;
	}
	defer response.Body.Close();

	responseBody, err := ioutil.ReadAll(response.Body);
	if err != nil {
		fmt.Printf("Error reading body response from Google Vision API: %v\n", err);
		responseMessage := "Could not handle request. Please try again later.";
		helpers.SendErrorResponse(responseWriter, http.StatusBadRequest, responseMessage);
		return;
	}

	result, err := parseResponseBodyToSuccessResponseFormat(responseBody);
	if err != nil {
		fmt.Printf("Error parsing JSON response from Google Vision API: %v\n", err);
		responseMessage := "Could not handle request. Please try again later.";
		helpers.SendErrorResponse(responseWriter, http.StatusBadRequest, responseMessage);
		return;
	}

	if resultContainsIdentificationLabelsForImage(result) {
		errorMessage := determineResponseErrorMessageFromVisionResponse(responseBody);
		helpers.SendErrorResponse(responseWriter, http.StatusBadRequest, errorMessage);
	}

	sendResponseIfImageOfDogOrNot(responseWriter, result);
}

func verifyImageUrl(request *http.Request) (string, error) {
	imageURL := request.URL.Query().Get("imageUrl");
	if isImage(imageURL) {
		return "", errors.New("invalid url format");
	}
	return imageURL, nil;
}

func isImage(imageURL string) bool {
	r, _ := regexp.Compile(`/^https?:\/\/.+\.(jpg|jpeg|png|webp|avif|gif|svg)$/`);
	return r.MatchString(imageURL);
}

func parseResponseBodyToSuccessResponseFormat(responseBody []byte) (helpers.VisionAPISuccessResponseFormat, error) {
	var result helpers.VisionAPISuccessResponseFormat;
	err := json.Unmarshal(responseBody, &result);
	if (err != nil) {
		return result, err;
	}
	return result, nil;
}

func getResultFromVisionApiWithURL(imageURL string, googleApiKey string) (*http.Response, error) {
	requestBody := []byte(`{
		"requests": [
			{
				"image": {
				"source": {
					"imageUri": "` + imageURL + `"
				}
				},
				"features": [
					{
					"type": "LABEL_DETECTION"
					}
				]
			}
		]
	}`);
	requestBodyBytes := bytes.NewBuffer(requestBody);
	response, err := retryablehttp.Post("https://vision.googleapis.com/v1/images:annotate?key="+googleApiKey, "application/json", requestBodyBytes);
	return response, err;
}

func resultContainsIdentificationLabelsForImage(result helpers.VisionAPISuccessResponseFormat) bool {
	return len(result.Responses[0].LabelAnnotations) == 0;
}

func determineResponseErrorMessageFromVisionResponse(responseBody []byte) string {
	const notAnImageResponseMessage = "Bad image data.";
	const inaccessibleURLResponseMessage = "The URL does not appear to be accessible by us. Please double check or download the content and pass it in.";

	var unsuccessfulResult helpers.VisionApiErrorResponseFormat;
	json.Unmarshal(responseBody, &unsuccessfulResult);
	visionApiErrorMessage := unsuccessfulResult.Responses[0].Error.Message;
	if visionApiErrorMessage == inaccessibleURLResponseMessage {
		return "URL is invalid or not accessible.";
	} else if visionApiErrorMessage == notAnImageResponseMessage {
		return "URL does not contain raw image data.";
	} else {
		fmt.Printf("Unknown error response from google: {%s}", visionApiErrorMessage);
		return "Unknown error occurred. Please check your data.";
	}
}

func sendResponseIfImageOfDogOrNot(responseWriter http.ResponseWriter, result helpers.VisionAPISuccessResponseFormat) {
	isDog := helpers.IsImageOfDogFromVisionData(result);
	helpers.SendIdentificationResponse(responseWriter, isDog);
}
