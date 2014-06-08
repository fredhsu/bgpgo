package bgpgo

import (
	//"fmt"
	"testing"
)

func TestIdle(t *testing.T) {
    in := input{"bgpstart"}
    newState, stateName := idle(&in)
    if stateName != "Connect" && newState != nil {
        t.Errorf("Wrong state")
    }
	}

