package main

import (
    "log"
    zmq "github.com/pebbe/zmq4"
)

func main() {
    log.Print("Starting")
    server, _ := zmq.NewSocket(zmq.REP)
    defer server.Close()
    server.Bind("tcp://*:5432")

    for {
        msg, _ := server.Recv(0)
        log.Print("Recv", msg)
        reply := "World"
        server.Send(reply, 0)
        log.Print("Sent", reply)
    }
}
