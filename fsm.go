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
    OPEN_CONFIRM
    ESTABLISHED
)

type Event int
const (
    ManualStart Event = iota
    ManualStop
    AutomaticStart
    ManualStart_with_PassiveTcpEstablishment
    AutomaticStart_with_PassiveTcpEstablishment
    AutomaticStart_with_DampPeerOscillations
    AutomaticStart_with_DampPeerOscillations_and_PassiveTcpEstablishment
    AutomaticStop
    ConnectRetryTimer_Expires
    HoldTimer_Expires
    KeepaliveTimer_Expires
    DelayOpenTimer_Expires
    IdleHoldTimer_Expires
    TcpConnection_Valid
    Tcp_CR_Invalid
    Tcp_CR_Acked
    TcpConnectionConfirmed
    TcpConnectionFails
    BGPOpen
    BGPOpen_with_DelayOpenTimer_running
    BGPHeaderErr
    BGPOpenMsgErr
    OpenCollisionDump
    NotifMsgVerErr
    NotifMsg
    KeepAliveMsg
    UpdateMsg
    UpdateMsgErr
)


type stateFn func(*input) (stateFn, State)

type input struct {
    event string
}

func idle(in *input) (stateFn, State) {
    if in.event == "bgpstart" {
        return connect, CONNECT
    } 
    return idle, IDLE
}

func connect(in *input) (stateFn, State) {
    switch in.event {
    case "bgpstart", "connectRetryExpired":
        return connect, CONNECT
    case "transportConnectionOpen":
        return openSent, OPEN_SENT
    case "transportOpenFailed":
        return active, ACTIVE
    }
    return idle, IDLE
}

func active(in *input) (stateFn, State) {
    switch in.event {
    case "bgpstart", "transportOpenFailed":
        return active, ACTIVE
    case "connectRetryExpired":
        return connect, CONNECT
    case "transportConnectionOpen":
        return openSent, OPEN_SENT
    }
    return idle, IDLE
}

func openSent(in *input) (stateFn, State) {
    switch in.event {
    case "bgpstart":
        return openSent, OPEN_SENT
        case "receiveOpen":
            return openConfirm, OPEN_CONFIRM
        case "connectionClosed":
            return active, ACTIVE
    }
    return idle, IDLE
}

func openConfirm(in *input) (stateFn, State) {
    switch in.event {
    case "bgpstart", "keepaliveTimerExpired":
        return openConfirm, OPEN_CONFIRM
        case "receiveKeepalive":
            return established, ESTABLISHED
    }
    return idle, IDLE
}

func established(in *input) (stateFn, State) {
    switch in.event {
        case "bgpStart", "keepaliveTimerExpired", "receiveKeepalive", "receiveUpdate":
            return established, ESTABLISHED
    }
    return idle, IDLE
}


func run(j chan *input) {
    startState := idle
    stateName := IDLE
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
