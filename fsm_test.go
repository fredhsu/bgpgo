package bgpgo

import (
	//"fmt"
	"testing"
)

func TestIdle(t *testing.T) {
	//in := input{START}
	for i := 0; i < 13; i++ {
		in := input{Event(i), nil}
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

func TestActive(t *testing.T) {
	for i := 0; i < 13; i++ {
		in := input{Event(i), nil}
		_, stateName := active(&in)
		switch Event(i) {
		case START, TRANSPORT_FAILED:
			if stateName != ACTIVE {
				t.Errorf("Wrong State: ", stateName)
			}
			continue
		case CONNECT_RETRY_EXPIRED:
			if stateName != CONNECT {
				t.Errorf("Wrong State")
			}
			continue
		case TRANSPORT_OPEN:
			if stateName != OPEN_SENT {
				t.Errorf("Wrong State")
			}
			continue
		}
		if stateName != IDLE {
			t.Errorf("Wrong State")
		}
	}
}

func TestConnect(t *testing.T) {
	for i := 0; i < 13; i++ {
		in := input{Event(i), nil}
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

func TestOpenSent(t *testing.T) {
	for i := 0; i < 13; i++ {
		in := input{Event(i), nil}
		_, stateName := openSent(&in)
		switch Event(i) {
		case START:
			if stateName != OPEN_SENT {
				t.Errorf("Wrong State")
			}
			continue
		case OPEN_RECV:
			if stateName != OPEN_CONFIRM {
				t.Errorf("Wrong State")
			}
			continue
		case TRANSPORT_CLOSED:
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

func TestOpenConfirm(t *testing.T) {
	for i := 0; i < 13; i++ {
		in := input{Event(i), nil}
		_, stateName := openConfirm(&in)
		switch Event(i) {
		case START, KEEPALIVE_EXPIRED:
			if stateName != OPEN_CONFIRM {
				t.Errorf("Wrong State")
			}
			continue
		case KEEPALIVE_RECV:
			if stateName != ESTABLISHED {
				t.Errorf("Wrong State")
			}
			continue
			if stateName != IDLE {
				t.Errorf("Wrong State")
			}
		}
	}
}

func TestEstablished(t *testing.T) {
	for i := 0; i < 13; i++ {
		in := input{Event(i), nil}
		_, stateName := established(&in)
		switch Event(i) {
		case START, KEEPALIVE_EXPIRED, KEEPALIVE_RECV, UPDATE_RECV:
			if stateName != ESTABLISHED {
				t.Errorf("Wrong State")
			}
			continue
			if stateName != IDLE {
				t.Errorf("Wrong State")
			}
		}
	}
}
