// BGP Finite State Machine
// Inspired by: http://rspace.googlecode.com/hg/slide/lex.html#slide-32
// http://play.golang.org/p/Sj5W3MBwg2

// TODO: Create Enum for the inputs
// Summary of states
// http://www.freesoft.org/CIE/RFC/1771/52.htm

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
    START Event = iota
    STOP
    TRANSPORT_OPEN
    TRANSPORT_CLOSED
    TRANSPORT_FAILED
    TRANSPORT_ERROR
    CONNECT_RETRY_EXPIRED
    HOLD_EXPIRED
    KEEPALIVE_EXPIRED
    OPEN_RECV
    KEEPALIVE_RECV
    UPDATE_RECV
    NOTIFICATION_RECV
)


type stateFn func(*input) (stateFn, State)

type input struct {
    event Event
}

func idle(in *input) (stateFn, State) {
    if in.event == START {
        return connect, CONNECT
    } 
    return idle, IDLE
}

func connect(in *input) (stateFn, State) {
    switch in.event {
    case START, CONNECT_RETRY_EXPIRED:
        return connect, CONNECT
    case TRANSPORT_OPEN:
        return openSent, OPEN_SENT
    case TRANSPORT_FAILED:
        return active, ACTIVE
    }
    return idle, IDLE
}

func active(in *input) (stateFn, State) {
    switch in.event {
    case START, TRANSPORT_FAILED:
        return active, ACTIVE
    case CONNECT_RETRY_EXPIRED:
        return connect, CONNECT
    case TRANSPORT_OPEN:
        return openSent, OPEN_SENT
    }
    return idle, IDLE
}

func openSent(in *input) (stateFn, State) {
    switch in.event {
    case START:
        return openSent, OPEN_SENT
        case OPEN_RECV:
            return openConfirm, OPEN_CONFIRM
        case TRANSPORT_CLOSED:
            return active, ACTIVE
    }
    return idle, IDLE
}

func openConfirm(in *input) (stateFn, State) {
    switch in.event {
    case START, KEEPALIVE_EXPIRED:
        return openConfirm, OPEN_CONFIRM
        case KEEPALIVE_RECV:
            return established, ESTABLISHED
    }
    return idle, IDLE
}

func established(in *input) (stateFn, State) {
    switch in.event {
        case START, KEEPALIVE_EXPIRED, KEEPALIVE_RECV, UPDATE_RECV:
            return established, ESTABLISHED
    }
    return idle, IDLE
}


func runfsm(j chan *input) {
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
    in := input{START}
    go runfsm(ci)
    ci <- &in
    ci <- &input{TRANSPORT_OPEN}
    ci <- &input{OPEN_RECV}
    ci <- &input{KEEPALIVE_RECV}
    time.Sleep(100 * time.Second)
}
*/

/*
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
*/
