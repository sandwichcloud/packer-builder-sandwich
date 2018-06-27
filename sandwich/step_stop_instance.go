package sandwich

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/api/client"
)

type StepStopInstance struct {
}

func (s *StepStopInstance) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)
	instance := state.Get("instance").(*api.Instance)

	ui.Say("Stopping instance...")

	instanceClient := config.sandwichClient.Instance(config.ProjectName)
	err := instanceClient.ActionStop(instance.Name, false, 60)
	if err != nil {
		err := fmt.Errorf("Error stopping instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say("Waiting for instance to stop...")
	stateConf := &StateChangeConf{
		Pending:   []string{"POWERED_ON"},
		Target:    []string{"POWERED_OFF"},
		Refresh:   InstanceRefreshPowerFunc(instanceClient, instance.Name),
		StepState: state,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		err := fmt.Errorf("Error waiting for instance (%s) to stop: %s", instance.Name, err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say("Instance has stopped.")

	return multistep.ActionContinue
}

func (s *StepStopInstance) Cleanup(state multistep.StateBag) {
	// No cleanup
}

func InstanceRefreshPowerFunc(instanceClient client.InstanceClientInterface, instanceName string) func() (result interface{}, state string, err error) {
	return func() (result interface{}, state string, err error) {
		instance, err := instanceClient.Get(instanceName)
		if err != nil {
			if apiError, ok := err.(api.APIErrorInterface); ok {
				if apiError.IsNotFound() {
					return instance, "Deleted", nil
				}
			}
			return nil, "", err
		}
		return instance, instance.PowerState, nil
	}
}
