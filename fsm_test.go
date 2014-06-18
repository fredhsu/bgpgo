package bgpgo

import (
	//"fmt"
	"testing"
)

func TestIdle(t *testing.T) {
    //in := input{START}
    for i := 0; i < 13; i++ {
        in := input{Event(i)}
        _, stateName := idle(&in)
        if Event(i) == START {
            if stateName != CONNECT { 
                t.Errorf("Wrong state should be CONNECT")
            }
        } else {
            if stateName != IDLE {
                t.Errorf("Wrong state ", stateName, " should be IDLE")
            }
        }
    }
}

func TestConnect(t *testing.T) {
    for i := 0; i < 13; i++ {
        in := input{Event(i)}
        _, stateName := connect(&in)
        switch Event(i) {
        case START, CONNECT_RETRY_EXPIRED:
            if stateName != CONNECT {
                t.Errorf("Wrong State")
            }
            continue
        case TRANSPORT_OPEN:
            if stateName != OPEN_SENT {
                t.Errorf("Wrong State")
            }
            continue
        case TRANSPORT_FAILED:
            if stateName != ACTIVE {
                t.Errorf("Wrong State")
            }
            continue
        }
        if stateName != IDLE {
            t.Errorf("Wrong State")
        }
    }
}

