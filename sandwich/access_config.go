package sandwich

import (
	"github.com/hashicorp/packer/template/interpolate"
	"github.com/sandwichcloud/deli-cli/api/client"
	"github.com/sandwichcloud/deli-cli/utils"
)

type AccessConfig struct {
	APIServer  string `mapstructure:"api_server"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	AuthMethod string `mapstructure:"auth_method"`
	ProjectID  string `mapstructure:"project_id"`
	RegionID   string `mapstructure:"region_id"`

	sandwichClient client.ClientInterface
}

func (c *AccessConfig) Prepare(ctx *interpolate.Context) []error {

	c.sandwichClient = &client.SandwichClient{
		APIServer: &c.APIServer,
	}

	apiDiscover, err := c.sandwichClient.Auth().DiscoverAuth()
	if err != nil {
		return []error{err}
	}
	if c.AuthMethod == "" {
		c.AuthMethod = *apiDiscover.Default
	}

	token, err := utils.Login(c.sandwichClient.Auth(), c.Username, c.Password, c.AuthMethod, false)
	if err != nil {
		return []error{err}
	}

	c.sandwichClient.SetToken(token)
	if c.ProjectID != "" {
		project, err := c.sandwichClient.Project().Get(c.ProjectID)
		if err != nil {
			return []error{err}
		}
		token, err = c.sandwichClient.Auth().ScopeToken(project)
		if err != nil {
			return []error{err}
		}
		c.sandwichClient.SetToken(token)
	}

	return nil
}
