package bgpgo

import (
    "bytes"
    "net"
    "os"
    "encoding/binary"
    "fmt"
    "time"
)

const (
    RECV_BUF_LEN uint16 = 1024
    MSG_HEADER_LEN uint16 = 19
    OPEN_HEADER_LEN uint16 = 10
)


type MessageHeader struct {
    Marker [16]byte
    Length uint16
    Type uint8
}

func NewMessageHeader(l uint16, t uint8) *MessageHeader {
    marker := [16]byte{255, 255, 255,255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
    return &MessageHeader{Marker: marker, Length: l, Type: t}
}

// Keepalive is just the header

type OpenMessage struct {
    Version uint8
    MyAS uint16
    HoldTime uint16
    BGPId uint32
    OptParamLen uint8
    OptParms []Parameter
}

type Parameter struct {
    Type uint8
    Length uint8
    Value []byte
}

type UpdateMessage struct {
    WithdrawnLength uint16
    WithdrawnRoutes []byte
    TotalPathAttr uint16
    PathAttr []byte
    Nlri []byte
}


type NotificationMessage struct {
    ErrorCode uint8
    ErrorSubcode uint8
    Data []byte
}


func RecvMsg(conn net.Conn) *bytes.Reader { 
    b := make([]byte, RECV_BUF_LEN)
    _, err := conn.Read(b)
    buf := bytes.NewReader(b)
    if err != nil {
        println("Error reading:", err.Error())
        return nil
    }
    return buf
}

func OpenTransport() net.Conn {
    l, err := net.Listen("tcp", "0.0.0.0:179")
    if err != nil {
        println("error listening to bgp port", err.Error())
        os.Exit(1)
    }
    defer l.Close()
    conn, err := l.Accept()
    if err != nil {
        println("Error accepting connection:", err.Error())
        return nil
    }
    return conn
}

//func BgpSvr(conn net.Conn, fsmch chan *input) {
func BgpSvr() {
    fsmch := make(chan *input)
    go runfsm(fsmch)
    fsmch <- &input{START}
    conn := OpenTransport()
    fsmch <- &input{TRANSPORT_OPEN}

    msghdr := MessageHeader{}
    buf := RecvMsg(conn)
    binary.Read(buf, binary.BigEndian, &msghdr)
    fmt.Println("first msghdr:", msghdr)

    //Determine msg type
    fsmch <- &input{OPEN_RECV}
    openMsg := HandleOpenMessage(buf)
    fmt.Println("Open Message Received:")
    fmt.Println(openMsg)
    //send open back
    SendOpen(conn)
    // State: OpenSent
    fmt.Println("OpenSent")
    // Listen for Open Message From Peer
    // State: OpenConfirm
    // Listen for Keepalive
    buf = RecvMsg(conn)
    binary.Read(buf, binary.BigEndian, &msghdr)
    switch msghdr.Type {
    case 1:
        fsmch <- &input{OPEN_RECV}
        fmt.Println("open")
    case 2:
        fsmch <- &input{UPDATE_RECV}
        fmt.Println( "update")
    case 3:
        fsmch <- &input{NOTIFICATION_RECV}
        fmt.Println( "notification")
    case 4:
        fsmch <- &input{KEEPALIVE_RECV}
        fmt.Println( "keepalive")
    case 5:
        fsmch <- &input{OPEN_RECV}
        fmt.Println( "route-refresh")
    }
    for {
        SendKeepalive(conn)
        time.Sleep(10 * time.Second)
    }
    conn.Close()
}

func HandleOpenMessage(r *bytes.Reader) OpenMessage{
    openMsg := OpenMessage{}
    binary.Read(r, binary.BigEndian, &openMsg.Version)
    binary.Read(r, binary.BigEndian, &openMsg.MyAS)
    binary.Read(r, binary.BigEndian, &openMsg.HoldTime)
    binary.Read(r, binary.BigEndian, &openMsg.BGPId)
    binary.Read(r, binary.BigEndian, &openMsg.OptParamLen)
    param := Parameter{}
    remain := openMsg.OptParamLen
    binary.Read(r, binary.BigEndian, &param.Type)
    binary.Read(r, binary.BigEndian, &param.Length)
    param.Value = make([]byte, param.Length)
    binary.Read(r, binary.BigEndian, &param.Value)
    remain = remain - param.Length - 2
    fmt.Println(param, " buf: ", remain)
    return openMsg
}

func SendKeepalive(conn net.Conn) {
    fmt.Println("Sending Keepalive")
    writebuf := new(bytes.Buffer)
    msglen := MSG_HEADER_LEN
    msghdr := *NewMessageHeader(msglen, 4)
    err := binary.Write(writebuf, binary.BigEndian, &msghdr)
    if err != nil {
        fmt.Println("error", err)
    }
    fmt.Println(writebuf.Bytes())
    fmt.Println(len(writebuf.Bytes()))
    conn.Write(writebuf.Bytes())
}

func SendOpen(conn net.Conn) {
    writebuf := new(bytes.Buffer)
    msglen := MSG_HEADER_LEN + OPEN_HEADER_LEN + 0
    fmt.Println("Msg len: ", msglen)
    msghdr := *NewMessageHeader(msglen, 1)
    //as := 2
    //hold := 180
    //id := 100
    
    openMsg := OpenMessage{4, 2, 180, 100, 0, []Parameter{}}
    //openMsg := OpenMessage{4, as, hold, id, 0, []Parameter{}}
    err := binary.Write(writebuf, binary.BigEndian, &msghdr)
    err = binary.Write(writebuf, binary.BigEndian, &openMsg.Version)
    err = binary.Write(writebuf, binary.BigEndian, &openMsg.MyAS)
    err = binary.Write(writebuf, binary.BigEndian, &openMsg.HoldTime)
    err = binary.Write(writebuf, binary.BigEndian, &openMsg.BGPId)
    err = binary.Write(writebuf, binary.BigEndian, &openMsg.OptParamLen)
    if err != nil {
        fmt.Println("error", err)
    }
    fmt.Println(writebuf.Bytes())
    fmt.Println(len(writebuf.Bytes()))
    conn.Write(writebuf.Bytes())
}
/*
func main() {
    println("Starting router")

    l, err := net.Listen("tcp", "0.0.0.0:179")
    if err != nil {
        println("error listening to bgp port", err.Error())
        os.Exit(1)
    }
    defer l.Close()
    //for {
        conn, err := l.Accept()
        if err != nil {
            println("Error accepting connection:", err.Error())
            return
        }
        BgpSvr(conn)
    //}
}
*/
