package main

import (
    "fmt"
    zmq "github.com/pebbe/zmq4"
)

func main() {
    fmt.Println("Connecting");
    req, _ := zmq.NewSocket(zmq.REQ)
    defer req.Close()

    req.Connect("tcp://localhost:5432")
    fmt.Println("Connected");

    msg := "Hello"
    fmt.Println("sending ", msg)
    req.Send(msg, 0)
    fmt.Println("Sent ", msg)

    reply, _ := req.Recv(0)
    fmt.Println("Recv ", reply)
}
