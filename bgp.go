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
    /*
    openMsg := OpenMessage{}
    binary.Read(buf, binary.BigEndian, &openMsg.Version)
    binary.Read(buf, binary.BigEndian, &openMsg.MyAS)
    binary.Read(buf, binary.BigEndian, &openMsg.HoldTime)
    binary.Read(buf, binary.BigEndian, &openMsg.BGPId)
    binary.Read(buf, binary.BigEndian, &openMsg.OptParamLen)
    fmt.Println(openMsg)
    openMsg := HandleOpenMessage(buf)
    param := Parameter{}
    remain := openMsg.OptParamLen
    binary.Read(buf, binary.BigEndian, &param.Type)
    binary.Read(buf, binary.BigEndian, &param.Length)
    param.Value = make([]byte, param.Length)
    binary.Read(buf, binary.BigEndian, &param.Value)
    remain = remain - param.Length - 2
    fmt.Println(param, " buf: ", remain)
    if remain > 0 {
        fmt.Println("More stuff!")
    } else {
        fmt.Println("Finished!")
        conn.Close()
    }
    */
    openMsg := HandleOpenMessage(buf)
    fmt.Println(openMsg)
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