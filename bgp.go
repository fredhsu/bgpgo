package main

import (
    "bytes"
    "net"
    "os"
    "encoding/binary"
    "fmt"
)

const (
    RECV_BUF_LEN = 1024
)

type MessageHeader struct {
    Marker [16]byte
    Length uint16
    Type uint8
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

type KeepaliveMessage struct {
}

type NotificationMessage struct {
    ErrorCode uint8
    ErrorSubcode uint8
    Data []byte
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

func BgpSvr(conn net.Conn) {
    msghdr := MessageHeader{}
    b := make([]byte, RECV_BUF_LEN)
    _, err := conn.Read(b)
    buf := bytes.NewReader(b)
    if err != nil {
        println("Error reading:", err.Error())
        return
    }
    binary.Read(buf, binary.BigEndian, &msghdr)
    fmt.Println(msghdr)
    //Determine msg type
    openMsg := HandleOpenMessage(buf)
    fmt.Println("Open Message:")
    fmt.Println(openMsg)
    //send open back
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

func SendOpen() {
    //b := [16]byte{255, 255, 255,255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
    //msghdr := MessageHeader{b, 43, 1}
    openMsg := OpenMessage{4, 1, 180, 100, 0, []Parameter{}}
    fmt.Println(openMsg)
    // Build the packet - msghdr + openMsg
}
