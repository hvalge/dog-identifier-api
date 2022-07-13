package helpers;

import (

);

type VisionAPISuccessResponseFormat struct {
	Responses []struct {
		LabelAnnotations []struct {
			Mid         string  `json:"mid"`
			Description string  `json:"description"`
			Score       float64 `json:"score"`
			Topicality  float64 `json:"topicality"`
		}
	}
}

type VisionApiErrorResponseFormat struct {
	Responses []struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
	}
}

func IsImageOfDogFromVisionData(result VisionAPISuccessResponseFormat) bool {
	for _, label := range result.Responses[0].LabelAnnotations {
		if label.Description == "Dog" && label.Score > 0.7 {
			return true;
		}
	}
	return false;
}
