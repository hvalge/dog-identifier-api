package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"github.com/hashicorp/go-retryablehttp"
)

func IsUrlOfDog(w http.ResponseWriter, r *http.Request) {
	googleApiKey := os.Getenv("IMAGE_IDENTIFIER_GOOGLE_API_KEY")
	w.Header().Set("Content-Type", "application/json")

	imageURL := r.URL.Query().Get("imageUrl")
	if !strings.HasPrefix(imageURL, "http://") && !strings.HasPrefix(imageURL, "https://") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": "URL is invalid."
			}`))
		return
	}

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
	}`)
	requestBodyBytes := bytes.NewBuffer(requestBody)
	response, err := retryablehttp.Post("https://vision.googleapis.com/v1/images:annotate?key="+googleApiKey, "application/json", requestBodyBytes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": "Failed to request Google Vision API. Please try again later."
			}`))
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": "Something went wrong with parsing response from Google Vision API."
			}`))
		fmt.Printf("Error reading body response from Google Vision API: %v\n", err)
		return
	}

	var successfulResult successfulVisionAPIResult
	err = json.Unmarshal(responseBody, &successfulResult)
	if (err != nil) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": "Got invalid JSON syntax from Google API."
			}`))
		fmt.Printf("Error parsing response from Google Vision API: %v\n", err)
	}

	var gotImageResult bool = verifyResultAndRespondOnFailure(successfulResult, responseBody, w)
	if !gotImageResult {
		return
	}

	isDog := false
	for _, label := range successfulResult.Responses[0].LabelAnnotations {
		if label.Description == "Dog" && label.Score > 0.7 {
			isDog = true
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if isDog {
		w.Write([]byte(`{
			"message": "Image contains a dog."
		}`))
	} else {
		w.Write([]byte(`{
			"message": "Image does not contain a dog."
		}`))
	}
}

type successfulVisionAPIResult struct {
	Responses []struct {
		LabelAnnotations []struct {
			Mid         string  `json:"mid"`
			Description string  `json:"description"`
			Score       float64 `json:"score"`
			Topicality  float64 `json:"topicality"`
		}
	}
}

type unsuccessfulVisionAPIResult struct {
	Responses []struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
	}
}

func verifyResultAndRespondOnFailure(successfulResult successfulVisionAPIResult, responseBody []byte, w http.ResponseWriter) bool {
	if len(successfulResult.Responses[0].LabelAnnotations) == 0 {
		const notAnImageResponseMessage = "Bad image data."
		const inaccessibleURLResponseMessage = "The URL does not appear to be accessible by us. Please double check or download the content and pass it in."

		w.WriteHeader(http.StatusBadRequest)
		
		var unsuccessfulResult unsuccessfulVisionAPIResult
		json.Unmarshal(responseBody, &unsuccessfulResult)
		if unsuccessfulResult.Responses[0].Error.Message == notAnImageResponseMessage {
			w.Write([]byte(`{
				"error": "URL does not contain raw image data."
		}`))
		} else if unsuccessfulResult.Responses[0].Error.Message == inaccessibleURLResponseMessage {
			w.Write([]byte(`{
				"error": "URL is invalid or not accessible."
		}`))
		} else {
			w.Write([]byte(`{
				"error": "Something went wrong."
				"body": ` + string(responseBody) + `
		}`))
		}

		return false
	} else {
		return true
	}
}
