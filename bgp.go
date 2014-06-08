package bgpgo

import (
    "bytes"
    "net"
    "os"
    "encoding/binary"
    "fmt"
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

// Keepalive is just the header
type KeepaliveMessage struct {
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


func BgpSvr(conn net.Conn) {
    msghdr := MessageHeader{}
    buf := RecvMsg(conn)
    binary.Read(buf, binary.BigEndian, &msghdr)
    fmt.Println(msghdr)
    //Determine msg type
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
    binary.Read(buf, binary.BigEndian, &msghdr)
    fmt.Println(msghdr)
    for {
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
