package instance

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

type InstanceClient struct {
	APIServer   *string
	HttpClient  *http.Client
	ProjectName string
}

func (client *InstanceClient) Create(name, imageName, regionName, zoneName, networkName, serviceAccountName, flavorName string, disk int, keypairNames []string, initialVolumes []api.InstanceInitialVolume, tags map[string]string, userData string) (*api.Instance, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type createBody struct {
		Name               string                      `json:"name"`
		ImageName          string                      `json:"image_name"`
		RegionName         string                      `json:"region_name"`
		ZoneName           string                      `json:"zone_name,omitempty"`
		ServiceAccountName string                      `json:"service_account_name,omitempty"`
		NetworkName        string                      `json:"network_name"`
		FlavorName         string                      `json:"flavor_name"`
		Disk               int                         `json:"disk,omitempty"`
		KeypairNames       []string                    `json:"keypair_names,omitempty"`
		InitialVolumes     []api.InstanceInitialVolume `json:"initial_volumes"`
		Tags               map[string]string           `json:"tags"`
		UserData           string                      `json:"user_data"`
	}

	body := createBody{
		Name:               name,
		ImageName:          imageName,
		RegionName:         regionName,
		ZoneName:           zoneName,
		NetworkName:        networkName,
		ServiceAccountName: serviceAccountName,
		FlavorName:         flavorName,
		Disk:               disk,
		KeypairNames:       keypairNames,
		InitialVolumes:     initialVolumes,
		Tags:               tags,
		UserData:           userData,
	}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, client.HttpClient, *client.APIServer+"/compute/v1/projects/"+client.ProjectName+"/instances", "application/json", bytes.NewBuffer(jsonBody))
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

	instance := &api.Instance{}
	json.Unmarshal(responseData, instance)
	return instance, nil
}

func (client *InstanceClient) Get(name string) (*api.Instance, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	response, err := ctxhttp.Get(ctx, client.HttpClient, *client.APIServer+"/compute/v1/projects/"+client.ProjectName+"/instances/"+name)
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

	instance := &api.Instance{}
	json.Unmarshal(responseData, instance)
	return instance, nil
}

func (client *InstanceClient) Delete(name string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	Url, err := url.Parse(*client.APIServer + "/compute/v1/projects/" + client.ProjectName + "/instances/" + name)
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

func (client *InstanceClient) List(imageName string, limit int, marker string) (*api.InstanceList, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	parameters := url.Values{}

	if len(imageName) > 0 {
		parameters.Add("image_name", imageName)
	}

	parameters.Add("limit", strconv.FormatInt(int64(limit), 10))

	if len(marker) > 0 {
		parameters.Add("marker", marker)
	}

	Url, err := url.Parse(*client.APIServer + "/compute/v1/projects/" + client.ProjectName + "/instances")
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

	instances := &api.InstanceList{}
	json.Unmarshal(responseData, instances)
	return instances, nil
}

func (client *InstanceClient) ActionStop(name string, hard bool, timeout int) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type stopBody struct {
		Hard    bool `json:"hard"`
		Timeout int  `json:"timeout"`
	}

	body := stopBody{Hard: hard, Timeout: timeout}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPut, *client.APIServer+fmt.Sprintf("/compute/v1/projects/"+client.ProjectName+"/instances/%s/action/stop", name), bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
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

	if response.StatusCode != http.StatusAccepted {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return err
		}
		return apiError
	}
	return nil
}

func (client *InstanceClient) ActionStart(name string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	req, err := http.NewRequest(http.MethodPut, *client.APIServer+fmt.Sprintf("/compute/v1/projects/"+client.ProjectName+"/instances/%s/action/start", name), nil)
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

	if response.StatusCode != http.StatusAccepted {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return err
		}
		return apiError
	}
	return nil
}

func (client *InstanceClient) ActionRestart(name string, hard bool, timeout int) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type restartBody struct {
		Hard    bool `json:"hard"`
		Timeout int  `json:"timeout"`
	}

	body := restartBody{Hard: hard, Timeout: timeout}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPut, *client.APIServer+fmt.Sprintf("/compute/v1/projects/"+client.ProjectName+"/instances/%s/action/restart", name), bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
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

	if response.StatusCode != http.StatusAccepted {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return err
		}
		return apiError
	}
	return nil
}

func (client *InstanceClient) ActionImage(instanceName string, imageName string) (*api.Image, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type createBody struct {
		Name string `json:"name"`
	}

	body := createBody{Name: imageName}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, client.HttpClient, *client.APIServer+fmt.Sprintf("/compute/v1/projects/"+client.ProjectName+"/instances/%s/action/image", instanceName), "application/json", bytes.NewBuffer(jsonBody))
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

	image := &api.Image{}
	json.Unmarshal(responseData, image)
	return image, nil
}
