package sandwich

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/api/client"
)

type StepRunInstance struct {
	Name            string
	FlavorName      string
	Disk            int
	SourceImageName string
	NetworkName     string
	UserData        string
	Tags            map[string]string

	instance *api.Instance
}

func (s *StepRunInstance) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {

	config := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	instanceClient := config.sandwichClient.Instance(config.ProjectName)

	ui.Say("Launching instance...")

	var keyPairNames []string
	keyPairName, hasKey := state.GetOk("keyPairName")
	if hasKey {
		keyPairNames = append(keyPairNames, keyPairName.(string))
	}

	if s.Tags == nil {
		s.Tags = make(map[string]string)
	}

	instance, err := instanceClient.Create(s.Name, s.SourceImageName, config.RegionName, "", s.NetworkName, "", s.FlavorName, s.Disk, keyPairNames, []api.InstanceInitialVolume{}, s.Tags, s.UserData)
	if err != nil {
		err := fmt.Errorf("Error creating source instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message(fmt.Sprintf("Instance Name: %s", instance.Name))
	ui.Say("Waiting for instance to become active...")

	stateConf := &StateChangeConf{
		Pending:   []string{"ToCreate", "Creating"},
		Target:    []string{"Created"},
		Refresh:   InstanceRefreshFunc(instanceClient, instance.Name),
		StepState: state,
	}
	latestInstance, err := stateConf.WaitForState()
	if err != nil {
		err := fmt.Errorf("Error waiting for instance (%s) to become ready: %s", instance.Name, err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say("Instance has become active!")
	s.instance = latestInstance.(*api.Instance)
	state.Put("instance", s.instance)

	return multistep.ActionContinue
}

func (s *StepRunInstance) Cleanup(state multistep.StateBag) {
	if s.instance == nil {
		return
	}

	config := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	instanceClient := config.sandwichClient.Instance(config.ProjectName)
	ui.Say(fmt.Sprintf("Terminating the source instance: %s ...", s.instance.Name))
	err := instanceClient.Delete(s.instance.Name)
	if err != nil {
		ui.Error(fmt.Sprintf("Error terminating instance (%s), may still be around: %s", s.instance.Name, err))
		return
	}

	stateConf := &StateChangeConf{
		Pending: []string{"ToDelete", "Deleting"},
		Target:  []string{"Deleted"},
		Refresh: InstanceRefreshFunc(instanceClient, s.instance.Name),
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		ui.Error(fmt.Sprintf("Error waiting for instance termination (%s), may still be around: %s", s.instance.Name, err))
		return
	}

}

func InstanceRefreshFunc(instanceClient client.InstanceClientInterface, instanceName string) func() (result interface{}, state string, err error) {
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
		return instance, instance.State, nil
	}
}
