package sandwich

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/api/client"
)

type StepImageInstance struct {
}

func (s *StepImageInstance) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)
	instance := state.Get("instance").(*api.Instance)

	ui.Say("Creating an image of the instance...")

	instanceClient := config.sandwichClient.Instance(config.ProjectName)
	imageClient := config.sandwichClient.Image(config.ProjectName)
	image, err := instanceClient.ActionImage(instance.Name, config.ImageName)
	if err != nil {
		err := fmt.Errorf("Error imaging instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message(fmt.Sprintf("Image Name: %s", image.Name))
	ui.Say("Waiting image to be created...")
	stateConf := &StateChangeConf{
		Pending:   []string{"ToCreate", "Creating"},
		Target:    []string{"Created"},
		Refresh:   ImageRefreshFunc(imageClient, image.Name),
		StepState: state,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		err := fmt.Errorf("Error waiting for image (%s) to create: %s", image.Name, err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("imageName", image.Name)
	ui.Say("Image has been created.")

	return multistep.ActionContinue
}

func (s *StepImageInstance) Cleanup(state multistep.StateBag) {
	// No cleanup
}

func ImageRefreshFunc(imageClient client.ImageClientInterface, imageName string) func() (result interface{}, state string, err error) {
	return func() (result interface{}, state string, err error) {
		image, err := imageClient.Get(imageName)
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
