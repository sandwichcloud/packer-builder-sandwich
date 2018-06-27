package region

import (
	"net/http"

	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"

	"net/url"
	"strconv"

	"fmt"

	"github.com/sandwichcloud/deli-cli/api"
	"golang.org/x/net/context/ctxhttp"
)

type RegionClient struct {
	APIServer  *string
	HttpClient *http.Client
}

func (regionClient *RegionClient) Create(name, datacenter, imageDatastore, imageFolder string) (*api.Region, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type createBody struct {
		Name           string `json:"name"`
		Datacenter     string `json:"datacenter"`
		ImageDatastore string `json:"image_datastore"`
		ImageFolder    string `json:"image_folder,omitempty"`
	}

	body := createBody{
		Name:           name,
		Datacenter:     datacenter,
		ImageDatastore: imageDatastore,
		ImageFolder:    imageFolder,
	}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, regionClient.HttpClient, *regionClient.APIServer+"/location/v1/regions", "application/json", bytes.NewBuffer(jsonBody))
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

	region := &api.Region{}
	json.Unmarshal(responseData, region)
	return region, nil
}

func (regionClient *RegionClient) Get(name string) (*api.Region, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	response, err := ctxhttp.Get(ctx, regionClient.HttpClient, *regionClient.APIServer+"/location/v1/regions/"+name)
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

	region := &api.Region{}
	json.Unmarshal(responseData, region)
	return region, nil
}

func (regionClient *RegionClient) Delete(name string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	Url, err := url.Parse(*regionClient.APIServer + "/location/v1/regions/" + name)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", Url.String(), nil)
	if err != nil {
		return err
	}
	response, err := ctxhttp.Do(ctx, regionClient.HttpClient, req)
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

func (regionClient *RegionClient) List(limit int, marker string) (*api.RegionList, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	parameters := url.Values{}

	parameters.Add("limit", strconv.FormatInt(int64(limit), 10))

	if len(marker) > 0 {
		parameters.Add("marker", marker)
	}

	Url, err := url.Parse(*regionClient.APIServer + "/location/v1/regions")
	if err != nil {
		return nil, err
	}
	Url.RawQuery = parameters.Encode()

	response, err := ctxhttp.Get(ctx, regionClient.HttpClient, Url.String())
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

	regions := &api.RegionList{}
	json.Unmarshal(responseData, regions)
	return regions, nil
}

func (regionClient *RegionClient) ActionSchedule(name string, schedulable bool) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type stopBody struct {
		Schedulable bool `json:"schedulable"`
	}

	body := stopBody{Schedulable: schedulable}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPut, *regionClient.APIServer+fmt.Sprintf("/location/v1/regions/%s/action/schedule", name), bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	response, err := ctxhttp.Do(ctx, regionClient.HttpClient, req)
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
