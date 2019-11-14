package main

import (
    "fmt"
    zmq "github.com/pebbe/zmq4"
    "github.com/golang/protobuf/proto"
    "tradfrid/proto"
)

func main() {
    fmt.Println("Connecting");
    client, _ := zmq.NewSocket(zmq.REQ)
    defer client.Close()

    client.Connect("tcp://localhost:5432")
    fmt.Println("Connected");

    req := &arith.Req {
        Op: arith.OP_ADD.Enum(),
        A: proto.Uint32(1),
        B: proto.Uint32(2),
    }

    req_data, _ := proto.Marshal(req)
    fmt.Println("sending ", req_data)
    client.SendBytes(req_data, 0)
    fmt.Println("Sent ", req_data)

    resp_data, _ := client.RecvBytes(0)
    fmt.Println("Recv ", resp_data)
    resp := &arith.Resp{}
    _ = proto.Unmarshal(resp_data, resp)
    fmt.Println(*resp.R)
}
