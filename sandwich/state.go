package sandwich

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mitchellh/multistep"
)

// StateRefreshFunc is a function type used for StateChangeConf that is
// responsible for refreshing the item being watched for a state change.
//
// It returns three results. `result` is any object that will be returned
// as the final object after waiting for state change. This allows you to
// return the final updated object, for example an EC2 instance after refreshing
// it.
//
// `state` is the latest state of that object. And `err` is any error that
// may have happened while refreshing the state.
type StateRefreshFunc func() (result interface{}, state string, err error)

// StateChangeConf is the configuration struct used for `WaitForState`.
type StateChangeConf struct {
	Pending   []string
	Refresh   StateRefreshFunc
	StepState multistep.StateBag
	Target    []string
}

// WaitForState watches an object and waits for it to achieve a certain
// state.
func (conf *StateChangeConf) WaitForState() (i interface{}, err error) {
	log.Printf("Waiting for state to become: %s", conf.Target)

	for {
		var currentProgress int
		var currentState string
		i, currentState, err = conf.Refresh()
		if err != nil {
			return
		}

		for _, t := range conf.Target {
			if currentState == t {
				return
			}
		}

		if conf.StepState != nil {
			if _, ok := conf.StepState.GetOk(multistep.StateCancelled); ok {
				return nil, errors.New("interrupted")
			}
		}

		found := false
		for _, allowed := range conf.Pending {
			if currentState == allowed {
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("unexpected state '%s', wanted target '%s'", currentState, conf.Target)
		}

		log.Printf("Waiting for state to become: %s currently %s (%d%%)", conf.Target, currentState, currentProgress)
		time.Sleep(2 * time.Second)
	}
}