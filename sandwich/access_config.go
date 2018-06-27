package sandwich

import (
	"errors"

	"github.com/hashicorp/packer/template/interpolate"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/api/client"
	"golang.org/x/oauth2"
)

type AccessConfig struct {
	APIServer   string `mapstructure:"api_server"`
	Token       string `mapstructure:"token"`
	ProjectName string `mapstructure:"project_name"`
	RegionName  string `mapstructure:"region_name"`

	sandwichClient client.ClientInterface
}

func (c *AccessConfig) Prepare(ctx *interpolate.Context) []error {

	if c.APIServer == "" {
		return []error{errors.New("API server is required")}
	}

	c.sandwichClient = &client.SandwichClient{
		APIServer: &c.APIServer,
	}

	if c.Token == "" {
		return []error{errors.New("Token is required")}
	}

	token := &oauth2.Token{
		AccessToken: c.Token,
		TokenType:   "Bearer",
	}

	c.sandwichClient.SetToken(token)
	if c.ProjectName != "" {
		_, err := c.sandwichClient.Project().Get(c.ProjectName)
		if err != nil {
			if apiError, ok := err.(api.APIErrorInterface); ok {
				if apiError.IsNotFound() {
					return []error{errors.New("Configured project does not exist.")}
				}
			}
			return []error{err}
		}
	} else {
		return []error{errors.New("A project is required")}
	}

	if c.RegionName != "" {
		// TODO: validate region
	} else {
		return []error{errors.New("A region is required")}
	}

	return nil
}
