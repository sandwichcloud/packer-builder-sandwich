package zone

import (
	"net/http"

	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"

	"fmt"
	"net/url"
	"strconv"

	"github.com/sandwichcloud/deli-cli/api"
	"golang.org/x/net/context/ctxhttp"
)

type ZoneClient struct {
	APIServer  *string
	HttpClient *http.Client
}

func (zoneClient *ZoneClient) Create(name, regionName, vmCluster, vmDatastore, vmFolder string, coreProvisionPercent, ramProvisionPercent int) (*api.Zone, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type createBody struct {
		Name                 string `json:"name"`
		RegionName           string `json:"region_name"`
		VMCluster            string `json:"vm_cluster"`
		VMDatastore          string `json:"vm_datastore"`
		VMFolder             string `json:"vm_folder,omitempty"`
		CoreProvisionPercent int    `json:"core_provision_percent"`
		RamProvisionPercent  int    `json:"ram_provision_percent"`
	}

	body := createBody{
		Name:                 name,
		RegionName:           regionName,
		VMCluster:            vmCluster,
		VMDatastore:          vmDatastore,
		CoreProvisionPercent: coreProvisionPercent,
		RamProvisionPercent:  ramProvisionPercent,
	}

	if len(vmFolder) > 0 {
		body.VMFolder = vmFolder
	}

	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, zoneClient.HttpClient, *zoneClient.APIServer+"/location/v1/zones", "application/json", bytes.NewBuffer(jsonBody))
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

	zone := &api.Zone{}
	json.Unmarshal(responseData, zone)
	return zone, nil
}

func (zoneClient *ZoneClient) Get(name string) (*api.Zone, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	response, err := ctxhttp.Get(ctx, zoneClient.HttpClient, *zoneClient.APIServer+"/location/v1/zones/"+name)
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

	zone := &api.Zone{}
	json.Unmarshal(responseData, zone)
	return zone, nil
}

func (zoneClient *ZoneClient) Delete(name string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	Url, err := url.Parse(*zoneClient.APIServer + "/location/v1/zones/" + name)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", Url.String(), nil)
	if err != nil {
		return err
	}
	response, err := ctxhttp.Do(ctx, zoneClient.HttpClient, req)
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

func (zoneClient *ZoneClient) List(regionName string, limit int, marker string) (*api.ZoneList, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	parameters := url.Values{}

	if len(regionName) > 0 {
		parameters.Add("region_name", regionName)
	}

	parameters.Add("limit", strconv.FormatInt(int64(limit), 10))

	if len(marker) > 0 {
		parameters.Add("marker", marker)
	}

	Url, err := url.Parse(*zoneClient.APIServer + "/location/v1/zones")
	if err != nil {
		return nil, err
	}
	Url.RawQuery = parameters.Encode()

	response, err := ctxhttp.Get(ctx, zoneClient.HttpClient, Url.String())
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

	zones := &api.ZoneList{}
	json.Unmarshal(responseData, zones)
	return zones, nil
}

func (zoneClient *ZoneClient) ActionSchedule(name string, schedulable bool) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type stopBody struct {
		Schedulable bool `json:"schedulable"`
	}

	body := stopBody{Schedulable: schedulable}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPut, *zoneClient.APIServer+fmt.Sprintf("/location/v1/zones/%s/action/schedule", name), bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	response, err := ctxhttp.Do(ctx, zoneClient.HttpClient, req)
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
