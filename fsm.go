// BGP Finite State Machine
// Inspired by: http://rspace.googlecode.com/hg/slide/lex.html#slide-32
// http://play.golang.org/p/Sj5W3MBwg2

// TODO: Create Enum for the states and inputs

package bgpgo

import (
    "fmt"
)

type State int

const (
    IDLE State = iota
    CONNECT
    OPEN_SENT
    ACTIVE
    OPEN_CONFIM
    ESTABLISHED
)

type stateFn func(*input) (stateFn, string)

type input struct {
    event string
}

func idle(in *input) (stateFn, string) {
    if in.event == "bgpstart" {
        return connect, "Connect"
    } 
    return idle, "Idle"
}

func connect(in *input) (stateFn, string) {
    switch in.event {
    case "bgpstart", "connectRetryExpired":
        return connect, "Connect"
    case "transportConnectionOpen":
        return openSent, "OpenSent"
    case "transportOpenFailed":
        return active, "Active"
    }
    return idle, "Idle"
}

func active(in *input) (stateFn, string) {
    switch in.event {
    case "bgpstart", "transportOpenFailed":
        return active, "Active"
    case "connectRetryExpired":
        return connect, "Connect"
    case "transportConnectionOpen":
        return openSent, "OpenSent"
    }
    return idle, "Idle"
}

func openSent(in *input) (stateFn, string) {
    switch in.event {
    case "bgpstart":
        return openSent, "OpenSent"
        case "receiveOpen":
            return openConfirm, "OpenConfirm"
        case "connectionClosed":
            return active, "Active"
    }
    return idle, "Idle"
}

func openConfirm(in *input) (stateFn, string) {
    switch in.event {
    case "bgpstart", "keepaliveTimerExpired":
        return openConfirm, "OpenConfirm"
        case "receiveKeepalive":
            return established, "Established"
    }
    return idle, "Idle"
}

func established(in *input) (stateFn, string) {
    switch in.event {
        case "bgpStart", "keepaliveTimerExpired", "receiveKeepalive", "receiveUpdate":
            return established, "Established"
    }
    return idle, "Idle"
}


func run(j chan *input) {
    startState := idle
    stateName := "Idle"
    for state := startState; state != nil; {
        in := <-j
        state, stateName = state(in)
        fmt.Println("Current State: ", stateName)
    }
}
/*
func main() {
    // Change them to return a value + statefn?
    ci := make(chan *input)
    in := input{"bgpstart"}
    go run(ci)
    ci <- &in
    ci <- &input{"transportConnectionOpen"}
    ci <- &input{"receiveOpen"}
    ci <- &input{"receiveKeepalive"}
    time.Sleep(100 * time.Second)
}
*/
