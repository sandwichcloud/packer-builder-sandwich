package sandwich

import (
	"log"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
	"github.com/mitchellh/multistep"
)

const BuilderId = "sandwich"

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	AccessConfig `mapstructure:",squash"`
	ImageConfig  `mapstructure:",squash"`
	RunConfig    `mapstructure:",squash"`

	ctx interpolate.Context
}

type Builder struct {
	config Config
	runner multistep.Runner
}

func (builder *Builder) Prepare(raws ...interface{}) ([]string, error) {
	err := config.Decode(&builder.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &builder.config.ctx,
	}, raws...)
	if err != nil {
		return nil, err
	}

	var errs *packer.MultiError
	errs = packer.MultiErrorAppend(errs, builder.config.AccessConfig.Prepare(&builder.config.ctx)...)
	errs = packer.MultiErrorAppend(errs, builder.config.ImageConfig.Prepare(&builder.config.ctx)...)
	errs = packer.MultiErrorAppend(errs, builder.config.RunConfig.Prepare(&builder.config.ctx)...)

	if errs != nil && len(errs.Errors) > 0 {
		return nil, errs
	}

	return nil, nil
}

func (builder *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	state := new(multistep.BasicStateBag)
	state.Put("config", builder.config)
	state.Put("hook", hook)
	state.Put("ui", ui)

	steps := []multistep.Step{
		&StepKeyPair{
			Debug:                builder.config.PackerDebug,
			SSHAgentAuth:         builder.config.RunConfig.Comm.SSHAgentAuth,
			TemporaryKeyPairName: builder.config.TemporaryKeyPairName,
			KeyPairID:            builder.config.SSHKeyPairID,
			PrivateKeyFile:       builder.config.RunConfig.Comm.SSHPrivateKey,
		},
		&StepRunInstance{
			Name:          builder.config.ImageName,
			FlavorID:      builder.config.FlavorID,
			Disk:          builder.config.Disk,
			SourceImageID: builder.config.SourceImageID,
			NetworkID:     builder.config.NetworkID,
			UserData:      builder.config.UserData,
			Tags:          builder.config.Tags,
		},
		&communicator.StepConnect{
			Config: &builder.config.RunConfig.Comm,
			Host:   CommHost(builder.config.sandwichClient.NetworkPort()),
			SSHConfig: SSHConfig(
				builder.config.RunConfig.Comm.SSHAgentAuth,
				builder.config.RunConfig.Comm.SSHUsername,
				builder.config.RunConfig.Comm.SSHPassword,
			),
		},
		&common.StepProvision{},
		&StepStopInstance{},
		&StepImageInstance{},
	}

	builder.runner = common.NewRunner(steps, builder.config.PackerConfig, ui)
	builder.runner.Run(state)

	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	// If there are no images, then just return
	if _, ok := state.GetOk("image"); !ok {
		return nil, nil
	}

	artifact := &Artifact{
		ImageID:        state.Get("imageID").(string),
		BuilderIdValue: BuilderId,
		ImageClient:    builder.config.sandwichClient.Image(),
	}

	return artifact, nil
}

func (builder *Builder) Cancel() {
	if builder.runner != nil {
		log.Println("Cancelling the step runner...")
		builder.runner.Cancel()
	}
}
