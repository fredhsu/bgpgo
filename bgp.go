package bgpgo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

type MessageType int

const (
	RECV_BUF_LEN    uint16 = 8196
	MSG_HEADER_LEN  uint16 = 19
	OPEN_HEADER_LEN uint16 = 10
)

const (
 _ MessageType = iota
 OPEN
UPDATE
NOTIFICATION
KEEPALIVE
ROUTE_REFRESH
)

type MessageHeader struct {
	Marker [16]byte
	Length uint16
	Type   uint8
}

func NewMessageHeader(l uint16, t uint8) *MessageHeader {
	marker := [16]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	return &MessageHeader{Marker: marker, Length: l, Type: t}
}

// Keepalive is just the header

type OpenMessage struct {
	Version     uint8
	MyAS        uint16
	HoldTime    uint16
	BGPId       uint32
	OptParamLen uint8
	OptParms    []Parameter
}

type Parameter struct {
	Type   uint8
	Length uint8
	Value  []byte
}

type UpdateMessage struct {
	WithdrawnLength uint16
	WithdrawnRoutes []byte
	TotalPathAttr   uint16
	PathAttr        []byte
	Nlri            []byte
}

type NotificationMessage struct {
	ErrorCode    uint8
	ErrorSubcode uint8
	Data         []byte
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

func RecvMsg2(conn net.Conn) *bytes.Reader {
    // io.Copy(c, conn)
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
	//TODO:  Should start connection to peer as well as listen
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

// func HandleMessage(r *bytes.Reader) {
//     msghdr := MessageHeader{}
//     binary.Read(buf, binary.BigEndian, &msghdr)
//     switch msghdr.Type {
//     case 1:
//         fsmch <- &input{OPEN_RECV}
//         fmt.Println("open")
//         openMsg := HandleOpenMessage(buf)
//     case 2:
//         fsmch <- &input{UPDATE_RECV}
//         fmt.Println( "update")
//     case 3:
//         fsmch <- &input{NOTIFICATION_RECV}
//         fmt.Println( "notification")
//     case 4:
//         fsmch <- &input{KEEPALIVE_RECV}
//         fmt.Println( "keepalive")
//     case 5:
//         fsmch <- &input{OPEN_RECV}
//         fmt.Println( "route-refresh")
//     }
// }

func getMessageType(msghdr MessageHeader, buf *bytes.Reader) MessageType {
	fmt.Printf("received: ")
	switch msghdr.Type {
	case 1:
		fmt.Println("open")
        return OPEN
	case 2:
		fmt.Println("update")
        return UPDATE
	case 3:
		fmt.Println("notification")
        return NOTIFICATION
	case 4:
		fmt.Println("keepalive")
        return KEEPALIVE
	case 5:
		fmt.Println("route-refresh")
        return ROUTE_REFRESH
	}
    return 0
}

func BgpSvr() {
	fsmch := make(chan *input)
	go runfsm(fsmch)
	fsmch <- &input{START, nil}
	conn := OpenTransport()
	fsmch <- &input{TRANSPORT_OPEN, conn}

	// Wait for message
	msghdr := MessageHeader{}
	buf := RecvMsg(conn)
	binary.Read(buf, binary.BigEndian, &msghdr)
	fmt.Println("msghdr:", msghdr)

	//Determine msg type
	fsmch <- &input{OPEN_RECV, conn}
	openMsg := HandleOpenMessage(buf)
	fmt.Println("Open Message Received:")
	fmt.Println(openMsg)
	//send open back
	// SendOpen(conn)
	// State: OpenSent
	// State: OpenConfirm
	SendKeepalive(conn)
	// Listen for Keepalive
	SendUpdate(conn)

	buf = RecvMsg(conn)
	binary.Read(buf, binary.BigEndian, &msghdr)
	m := getMessageType(msghdr, buf)
    fmt.Println("one")
   if m == UPDATE {
        HandleUpdateMessage(buf)
    }

	buf = RecvMsg(conn)
	binary.Read(buf, binary.BigEndian, &msghdr)
	m = getMessageType(msghdr, buf)
    fmt.Println("two")
    if m == UPDATE {
        HandleUpdateMessage(buf)
    }

	//Send Update

	// switch msghdr.Type {
	// case 1:
	// 	fsmch <- &input{OPEN_RECV, conn}
	// 	fmt.Println("open")
	// case 2:
	// 	fsmch <- &input{UPDATE_RECV, conn}
	// 	HandleUpdateMessage(buf)
	// 	fmt.Println("update")
	// case 3:
	// 	fsmch <- &input{NOTIFICATION_RECV, conn}
	// 	fmt.Println("notification")
	// case 4:
	// 	fsmch <- &input{KEEPALIVE_RECV, conn}
	// 	fmt.Println("keepalive")
	// case 5:
	// 	fsmch <- &input{OPEN_RECV, conn}
	// 	fmt.Println("route-refresh")
	// }

	buf = RecvMsg(conn)
	binary.Read(buf, binary.BigEndian, &msghdr)
	m = getMessageType(msghdr, buf)
    fmt.Println("three")
   if m == UPDATE {
        HandleUpdateMessage(buf)
    }
	for {
		SendKeepalive(conn)
		time.Sleep(10 * time.Second)
		// May need to spawn this
		// Need to keep checking for messages
	}
	conn.Close()
}

func HandleOpenMessage(r *bytes.Reader) OpenMessage {
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

func HandleUpdateMessage(r *bytes.Reader) UpdateMessage {
	updateMsg := UpdateMessage{}
	binary.Read(r, binary.BigEndian, &updateMsg.WithdrawnLength)
	updateMsg.WithdrawnRoutes = make([]byte, updateMsg.WithdrawnLength)
	binary.Read(r, binary.BigEndian, &updateMsg.WithdrawnRoutes)
	binary.Read(r, binary.BigEndian, &updateMsg.TotalPathAttr)
	updateMsg.PathAttr = make([]byte, updateMsg.TotalPathAttr)
    binary.Read(r, binary.BigEndian, &updateMsg.PathAttr)
	fmt.Println(updateMsg)
	// TODO: Parse Update message
	return updateMsg
}

func HandleNotificationMessage(r *bytes.Reader) NotificationMessage {
	notificationMsg := NotificationMessage{}
	binary.Read(r, binary.BigEndian, &notificationMsg.ErrorCode)
	binary.Read(r, binary.BigEndian, &notificationMsg.ErrorSubcode)

	fmt.Println(notificationMsg)
	return notificationMsg
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
	if conn == nil {
		// connection may be nil when unit testing
		return
	}
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
	fmt.Println("Sending Open")
	fmt.Println(writebuf.Bytes())
	fmt.Println(len(writebuf.Bytes()))
	conn.Write(writebuf.Bytes())
}

func SendUpdate(conn net.Conn) {
	if conn == nil {
		// connection may be nil when unit testing
		return
	}
	writebuf := new(bytes.Buffer)
	msglen := MSG_HEADER_LEN + 4 + 0
	fmt.Println("Msg len: ", msglen)
	msghdr := *NewMessageHeader(msglen, 1)
	// WithdrawnLength uint16
	// WithdrawnRoutes []byte
	// TotalPathAttr   uint16
	// PathAttr        []byte
	// Nlri
	updateMsg := UpdateMessage{0, nil, 0, nil, nil}
	//openMsg := OpenMessage{4, as, hold, id, 0, []Parameter{}}
	err := binary.Write(writebuf, binary.BigEndian, &msghdr)
	err = binary.Write(writebuf, binary.BigEndian, &updateMsg.WithdrawnLength)
	err = binary.Write(writebuf, binary.BigEndian, &updateMsg.WithdrawnRoutes)
	err = binary.Write(writebuf, binary.BigEndian, &updateMsg.TotalPathAttr)
	err = binary.Write(writebuf, binary.BigEndian, &updateMsg.PathAttr)
	err = binary.Write(writebuf, binary.BigEndian, &updateMsg.Nlri)
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Println("Sending Update")
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
