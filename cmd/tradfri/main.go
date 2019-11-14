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

    req := &remote.Req{
        Type: remote.ReqType_GetDevices.Enum(),
    }
    req_data, _ := proto.Marshal(req)
    client.SendBytes(req_data, 0)

    resp_data, _ := client.RecvBytes(0)
    resp := &remote.GetDevicesResp{}
    _ = proto.Unmarshal(resp_data, resp)
    fmt.Println(*resp)
}
