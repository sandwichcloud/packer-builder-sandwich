package policy

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sandwichcloud/deli-cli/api"
	"golang.org/x/net/context/ctxhttp"
)

type PolicyClient struct {
	APIServer  *string
	HttpClient *http.Client
	Type       string
}

func (client *PolicyClient) Get() (*api.Policy, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	response, err := ctxhttp.Get(ctx, client.HttpClient, *client.APIServer+"/iam/v1/"+client.Type)
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

	policy := &api.Policy{}
	json.Unmarshal(responseData, policy)
	return policy, nil
}

func (client *PolicyClient) Set(policy api.Policy) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	jsonBody, _ := json.Marshal(policy)
	response, err := ctxhttp.Post(ctx, client.HttpClient, *client.APIServer+"/iam/v1/"+client.Type, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		if err == context.DeadlineExceeded {
			return api.ErrTimedOut
		}
		return err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return err
		}
		return apiError
	}

	return nil
}
