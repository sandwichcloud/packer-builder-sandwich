package sandwich

import (
	"fmt"

	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/api/client"
)

type StepImageInstance struct {
}

func (s *StepImageInstance) Run(state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)
	instance := state.Get("instance").(*api.Instance)

	ui.Say("Creating an image of the instance...")

	instanceClient := config.sandwichClient.Instance()
	imageClient := config.sandwichClient.Image()
	image, err := instanceClient.ActionImage(instance.ID.String(), config.ImageName, "PRIVATE")
	if err != nil {
		err := fmt.Errorf("Error imaging instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message(fmt.Sprintf("Image ID: %s", image.ID.String()))
	ui.Say("Waiting image to be created...")
	stateConf := &StateChangeConf{
		Pending:   []string{"ToCreate", "Creating"},
		Target:    []string{"Created"},
		Refresh:   ImageRefreshFunc(imageClient, image.ID.String()),
		StepState: state,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		err := fmt.Errorf("Error waiting for image (%s) to create: %s", image.ID.String(), err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("imageID", image.ID.String())

	ui.Say("Image has been created.")

	return multistep.ActionContinue
}

func (s *StepImageInstance) Cleanup(state multistep.StateBag) {
	// No cleanup
}

func ImageRefreshFunc(imageClient client.ImageClientInterface, imageID string) func() (result interface{}, state string, err error) {
	return func() (result interface{}, state string, err error) {
		image, err := imageClient.Get(imageID)
		if err != nil {
			if apiError, ok := err.(api.APIErrorInterface); ok {
				if apiError.IsNotFound() {
					return image, "Deleted", nil
				}
			}
			return nil, "", err
		}
		return image, image.State, nil
	}
}
