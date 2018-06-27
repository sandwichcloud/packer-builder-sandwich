package network

import (
	"net/http"

	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"

	"net"

	"net/url"

	"strconv"

	"github.com/sandwichcloud/deli-cli/api"
	"golang.org/x/net/context/ctxhttp"
)

type NetworkClient struct {
	APIServer  *string
	HttpClient *http.Client
}

func (client *NetworkClient) Create(name, regionName, portGroup, cidr string, gateway, poolStart, poolEnd net.IP, dnsServers []net.IP) (*api.Network, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type createBody struct {
		Name       string   `json:"name"`
		PortGroup  string   `json:"port_group"`
		CIDR       string   `json:"cidr"`
		Gateway    net.IP   `json:"gateway"`
		DNSServers []net.IP `json:"dns_servers"`
		PoolStart  net.IP   `json:"pool_start"`
		PoolEnd    net.IP   `json:"pool_end"`
		RegionName string   `json:"region_name"`
	}

	body := createBody{
		Name:       name,
		PortGroup:  portGroup,
		CIDR:       cidr,
		Gateway:    gateway,
		DNSServers: dnsServers,
		PoolStart:  poolStart,
		PoolEnd:    poolEnd,
		RegionName: regionName,
	}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, client.HttpClient, *client.APIServer+"/compute/v1/networks", "application/json", bytes.NewBuffer(jsonBody))
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

	network := &api.Network{}
	json.Unmarshal(responseData, network)
	return network, nil
}

func (client *NetworkClient) Get(name string) (*api.Network, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	response, err := ctxhttp.Get(ctx, client.HttpClient, *client.APIServer+"/compute/v1/networks/"+name)
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

	network := &api.Network{}
	json.Unmarshal(responseData, network)
	return network, nil
}

func (client *NetworkClient) Delete(name string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	Url, err := url.Parse(*client.APIServer + "/compute/v1/networks/" + name)
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

func (client *NetworkClient) List(region_name string, limit int, marker string) (*api.NetworkList, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	parameters := url.Values{}

	if len(region_name) > 0 {
		parameters.Add("region_name", region_name)
	}

	parameters.Add("limit", strconv.FormatInt(int64(limit), 10))

	if len(marker) > 0 {
		parameters.Add("marker", marker)
	}

	Url, err := url.Parse(*client.APIServer + "/compute/v1/networks")
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

	networks := &api.NetworkList{}
	json.Unmarshal(responseData, networks)
	return networks, nil
}
