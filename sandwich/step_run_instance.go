package sandwich

import (
	"fmt"

	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/api/client"
)

type StepRunInstance struct {
	Name          string
	FlavorID      string
	Disk          int
	SourceImageID string
	NetworkID     string
	UserData      string
	Tags          map[string]string

	instance *api.Instance
}

func (s *StepRunInstance) Run(state multistep.StateBag) multistep.StepAction {

	config := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	instanceClient := config.sandwichClient.Instance()

	ui.Say("Launching instance...")

	var keyPairIDs []string
	keyPairID, hasKey := state.GetOk("keyPairID")
	if hasKey {
		keyPairIDs = append(keyPairIDs, keyPairID.(string))
	}

	if s.Tags == nil {
		s.Tags = make(map[string]string)
	}

	instance, err := instanceClient.Create(s.Name, s.SourceImageID, config.RegionID, "", s.NetworkID, "", s.FlavorID, s.Disk, keyPairIDs, []api.InstanceInitialVolume{}, s.Tags, s.UserData)
	if err != nil {
		err := fmt.Errorf("Error creating source instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message(fmt.Sprintf("Instance ID: %s", instance.ID.String()))
	ui.Say("Waiting for instance to become active...")

	stateConf := &StateChangeConf{
		Pending:   []string{"ToCreate", "Creating"},
		Target:    []string{"Created"},
		Refresh:   InstanceRefreshFunc(instanceClient, instance.ID.String()),
		StepState: state,
	}
	latestInstance, err := stateConf.WaitForState()
	if err != nil {
		err := fmt.Errorf("Error waiting for instance (%s) to become ready: %s", instance.ID.String(), err)
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

	instanceClient := config.sandwichClient.Instance()
	ui.Say(fmt.Sprintf("Terminating the source instance: %s ...", s.instance.ID.String()))
	err := instanceClient.Delete(s.instance.ID.String())
	if err != nil {
		ui.Error(fmt.Sprintf("Error terminating instance (%s), may still be around: %s", s.instance.ID.String(), err))
		return
	}

	stateConf := &StateChangeConf{
		Pending: []string{"ToDelete", "Deleting"},
		Target:  []string{"Deleted"},
		Refresh: InstanceRefreshFunc(instanceClient, s.instance.ID.String()),
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		ui.Error(fmt.Sprintf("Error waiting for instance termination (%s), may still be around: %s", s.instance.ID.String(), err))
		return
	}

}

func InstanceRefreshFunc(instanceClient client.InstanceClientInterface, instanceID string) func() (result interface{}, state string, err error) {
	return func() (result interface{}, state string, err error) {
		instance, err := instanceClient.Get(instanceID)
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
