package bgpgo

import (
	//"fmt"
	"testing"
)

func TestIdle(t *testing.T) {
    in := input{START}
    newState, stateName := idle(&in)
    if stateName != CONNECT && newState != nil {
        t.Errorf("Wrong state")
    }
	}

