package bgpgo

import (
	//"fmt"
	"testing"
)

func TestIdle(t *testing.T) {
    in := input{"bgpstart"}
    newState, stateName := idle(&in)
    if stateName != CONNECT && newState != nil {
        t.Errorf("Wrong state")
    }
	}

