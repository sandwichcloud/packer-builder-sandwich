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
	SSHKeyPairID         string              `mapstructure:"ssh_keypair_id"`
	TemporaryKeyPairName string              `mapstructure:"temporary_keypair_name"`

	SourceImageID string            `mapstructure:"source_image_id"`
	FlavorID      string            `mapstructure:"flavor_id"`
	NetworkID     string            `mapstructure:"network_id"`
	UserData      string            `mapstructure:"user_data"`
	Tags          map[string]string `mapstructure:"tags"`
}

func (c *RunConfig) Prepare(ctx *interpolate.Context) []error {

	// If we are not given an explicit ssh_keypair_id or
	// ssh_private_key_file, then create a temporary one, but only if the
	// temporary_key_pair_name has not been provided and we are not using
	// ssh_password.
	if c.SSHKeyPairID == "" && c.TemporaryKeyPairName == "" && c.Comm.SSHPrivateKey == "" && c.Comm.SSHPassword == "" {
		c.TemporaryKeyPairName = fmt.Sprintf("packer_%s", uuid.TimeOrderedUUID())
	}

	errs := c.Comm.Prepare(ctx)

	if c.SSHKeyPairID != "" {
		if c.Comm.Type == "winrm" && c.Comm.WinRMPassword == "" && c.Comm.SSHPrivateKey == "" {
			errs = append(errs, errors.New("A ssh_private_key_file must be provided to retrieve the winrm password when using ssh_keypair_id."))
		} else if c.Comm.SSHPrivateKey == "" && !c.Comm.SSHAgentAuth {
			errs = append(errs, errors.New("A ssh_private_key_file must be provided or ssh_agent_auth enabled when ssh_keypair_id is specified."))
		}
	}

	if c.SourceImageID == "" {
		errs = append(errs, errors.New("source_image_id must be specified"))
	}

	if c.NetworkID == "" {
		errs = append(errs, errors.New("network_id must be specified"))
	}

	//if c.FlavorID == "" {
	//	errs = append(errs, errors.New("flavor_id must be specified"))
	//}

	return errs
}
