package sandwich

import (
	"errors"
	"fmt"

	"github.com/hashicorp/packer/common/uuid"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/template/interpolate"
)

type RunConfig struct {
	Comm                 communicator.Config `mapstructure:",squash"`
	SSHKeyPairName       string              `mapstructure:"ssh_keypair_name"`
	TemporaryKeyPairName string              `mapstructure:"temporary_keypair_name"`

	SourceImageName string            `mapstructure:"source_image_name"`
	InstanceName    string            `mapstructure:"instance_name"`
	FlavorName      string            `mapstructure:"flavor_name"`
	Disk            int               `mapstructure:"disk"`
	NetworkName     string            `mapstructure:"network_name"`
	UserData        string            `mapstructure:"user_data"`
	Tags            map[string]string `mapstructure:"tags"`
}

func (c *RunConfig) Prepare(ctx *interpolate.Context) []error {

	// If we are not given an explicit ssh_keypair_name or
	// ssh_private_key_file, then create a temporary one, but only if the
	// temporary_key_pair_name has not been provided and we are not using
	// ssh_password.
	if c.SSHKeyPairName == "" && c.TemporaryKeyPairName == "" && c.Comm.SSHPrivateKey == "" && c.Comm.SSHPassword == "" {
		c.TemporaryKeyPairName = fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())
	}

	errs := c.Comm.Prepare(ctx)

	if c.SSHKeyPairName != "" {
		if c.Comm.Type == "winrm" && c.Comm.WinRMPassword == "" && c.Comm.SSHPrivateKey == "" {
			errs = append(errs, errors.New("A ssh_private_key_file must be provided to retrieve the winrm password when using ssh_keypair_name."))
		} else if c.Comm.SSHPrivateKey == "" && !c.Comm.SSHAgentAuth {
			errs = append(errs, errors.New("A ssh_private_key_file must be provided or ssh_agent_auth enabled when ssh_keypair_name is specified."))
		}
	}

	if c.SourceImageName == "" {
		errs = append(errs, errors.New("source_image_name must be specified"))
	}

	if c.NetworkName == "" {
		errs = append(errs, errors.New("network_name must be specified"))
	}

	if c.FlavorName == "" {
		errs = append(errs, errors.New("flavor_name must be specified"))
	}

	if c.InstanceName == "" {
		c.InstanceName = fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())
	}

	return errs
}
