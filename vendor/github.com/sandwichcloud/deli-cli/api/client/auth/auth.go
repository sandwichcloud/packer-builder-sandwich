package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"context"

	"github.com/sandwichcloud/deli-cli/api"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/oauth2"
)

type AuthClient struct {
	APIServer  *string
	HttpClient *http.Client
}

func (authClient *AuthClient) Login(username, password string) (*oauth2.Token, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	data := requestBody{Username: username, Password: password}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("Error parsing auth data into json")
	}

	response, err := ctxhttp.Post(ctx, http.DefaultClient, *authClient.APIServer+"/auth/v1/oauth/token", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, api.ErrTimedOut
		}
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return nil, err
		}
		return nil, apiError
	}

	token := &oauth2.Token{}
	json.Unmarshal(responseData, token)
	return token, nil
}
