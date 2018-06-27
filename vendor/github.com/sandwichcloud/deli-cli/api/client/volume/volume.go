package volume

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sandwichcloud/deli-cli/api"
	"golang.org/x/net/context/ctxhttp"
)

type VolumeClient struct {
	APIServer   *string
	HttpClient  *http.Client
	ProjectName string
}

func (client *VolumeClient) Create(name, zoneName string, size int) (*api.Volume, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type createBody struct {
		Name     string `json:"name"`
		ZoneName string `json:"zone_name"`
		Size     int    `json:"size"`
	}

	body := createBody{
		Name:     name,
		ZoneName: zoneName,
		Size:     size,
	}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, client.HttpClient, *client.APIServer+"/compute/v1/projects/"+client.ProjectName+"/volumes", "application/json", bytes.NewBuffer(jsonBody))
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

	volume := &api.Volume{}
	json.Unmarshal(responseData, volume)
	return volume, nil
}

func (client *VolumeClient) Get(name string) (*api.Volume, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	response, err := ctxhttp.Get(ctx, client.HttpClient, *client.APIServer+"/compute/v1/projects/"+client.ProjectName+"/volumes/"+name)
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

	volume := &api.Volume{}
	json.Unmarshal(responseData, volume)
	return volume, nil
}

func (client *VolumeClient) Delete(name string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	Url, err := url.Parse(*client.APIServer + "/compute/v1/projects/" + client.ProjectName + "/volumes/" + name)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", Url.String(), nil)
	if err != nil {
		return err
	}
	response, err := ctxhttp.Do(ctx, client.HttpClient, req)
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

func (client *VolumeClient) List(limit int, marker string) (*api.VolumeList, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	parameters := url.Values{}

	parameters.Add("limit", strconv.FormatInt(int64(limit), 10))

	if len(marker) > 0 {
		parameters.Add("marker", marker)
	}

	Url, err := url.Parse(*client.APIServer + "/compute/v1/projects/" + client.ProjectName + "/volumes")
	if err != nil {
		return nil, err
	}
	Url.RawQuery = parameters.Encode()

	response, err := ctxhttp.Get(ctx, client.HttpClient, Url.String())
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

	volumes := &api.VolumeList{}
	json.Unmarshal(responseData, volumes)
	return volumes, nil
}

func (client *VolumeClient) ActionAttach(name, instanceName string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type attachBody struct {
		InstanceName string `json:"instance_name"`
	}

	body := attachBody{InstanceName: instanceName}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, client.HttpClient, *client.APIServer+fmt.Sprintf("/compute/v1/projects/%s/volumes/%s/action/attach", client.ProjectName, name), "application/json", bytes.NewBuffer(jsonBody))
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

func (client *VolumeClient) ActionDetach(name string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	req, err := http.NewRequest(http.MethodPut, *client.APIServer+fmt.Sprintf("/compute/v1/projects/%s/volumes/%s/action/attach", client.ProjectName, name), nil)
	if err != nil {
		return err
	}
	response, err := ctxhttp.Do(ctx, client.HttpClient, req)
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

func (client *VolumeClient) ActionGrow(name string, newSize int) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type growBody struct {
		Size int `json:"size"`
	}

	body := growBody{Size: newSize}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, client.HttpClient, *client.APIServer+fmt.Sprintf("/compute/v1/projects/%s/volumes/%s/action/attach", client.ProjectName, name), "application/json", bytes.NewBuffer(jsonBody))
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

func (client *VolumeClient) ActionClone(name, newName string) (*api.Volume, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type cloneBody struct {
		Name string `json:"name"`
	}

	body := cloneBody{Name: newName}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, client.HttpClient, *client.APIServer+fmt.Sprintf("/compute/v1/projects/%s/volumes/%s/action/attach", client.ProjectName, name), "application/json", bytes.NewBuffer(jsonBody))
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

	volume := &api.Volume{}
	json.Unmarshal(responseData, volume)
	return volume, nil
}
