package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

var ErrTimedOut = errors.New("Timed out when communicating with the api server.")

type APIErrorInterface interface {
	error
	ToRawJSON() string
	IsNotFound() bool
}

type APIError struct {
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Errors     []dataError `json:"errors,omitempty"`
}

type dataError struct {
	Detail string `json:"detail"`
	Source source `json:"source"`
}

type source struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

func ParseErrors(statusCode int, errorData []byte) (APIError, error) {

	apiError := APIError{
		StatusCode: statusCode,
	}
	err := json.Unmarshal(errorData, &apiError)
	if err != nil {
		return apiError, err
	}

	return apiError, err
}

func (apiError APIError) IsNotFound() bool {
	if apiError.StatusCode == http.StatusNotFound {
		return true
	}
	return false
}

func (apiError APIError) ToRawJSON() string {
	errorData, _ := json.MarshalIndent(apiError, "", "  ")
	return string(errorData)
}

func (apiError APIError) Error() string {

	if apiError.Errors == nil {
		return apiError.Message
	}

	var errorStrings []string
	for _, dataError := range apiError.Errors {
		if dataError.Source.Pointer != "" {
			pointer := strings.Replace(dataError.Source.Pointer, "/data/attributes/", "", 1)
			errorStrings = append(errorStrings, pointer+": "+dataError.Detail)
		} else {
			errorStrings = append(errorStrings, dataError.Source.Parameter+": "+dataError.Detail)
		}

	}

	return strings.Join(errorStrings, ", ")
}
