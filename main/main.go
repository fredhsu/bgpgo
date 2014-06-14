package main

import (
    "net"
    "os"
    "github.com/fredhsu/bgpgo"
)

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
        bgpgo.BgpSvr(conn)
    //}
}
