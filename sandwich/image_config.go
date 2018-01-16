package sandwich

import (
	"fmt"

	"github.com/hashicorp/packer/template/interpolate"
)

type ImageConfig struct {
	ImageName    string   `mapstructure:"image_name"`
	ImagePublic  bool     `mapstructure:"image_public"`
	ImageMembers []string `mapstructure:"image_members"`
}

func (c *ImageConfig) Prepare(ctx *interpolate.Context) []error {

	errs := make([]error, 0)
	if c.ImageName == "" {
		errs = append(errs, fmt.Errorf("An image_name must be specified"))
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
