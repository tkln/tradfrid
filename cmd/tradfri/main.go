package main

import (
    "fmt"
    "flag"
    "strconv"
    zmq "github.com/pebbe/zmq4"
    "github.com/golang/protobuf/proto"
    "github.com/tkln/tradfrid/api"
)

func main() {
    flag.Parse()
    args := flag.Args()

    fmt.Println("Connecting");
    client, _ := zmq.NewSocket(zmq.REQ)
    defer client.Close()

    client.Connect("tcp://10.23.1.9:5432")
    //client.Connect("tcp://localhost:5432")
    fmt.Println("Connected");

    req := &remote.Req{
        Type: remote.ReqType_GetDevices.Enum(),
    }
    req_data, _ := proto.Marshal(req)
    client.SendBytes(req_data, 0)

    resp_data, _ := client.RecvBytes(0)
    resp := &remote.GetDevicesResp{}
    _ = proto.Unmarshal(resp_data, resp)
    devices := resp.Devices
    fmt.Println(devices)

    if len(args) > 1 {
        state := &remote.SetDeviceStateReq{}
        switch args[0] {
        case "on":
            state.Data = &remote.SetDeviceStateReq_Onoff{true}
        case "off":
            state.Data = &remote.SetDeviceStateReq_Onoff{false}
        case "level":
            if len(args) > 1 {
                val, _ := strconv.ParseFloat(args[1], 32)
                state.Data = &remote.SetDeviceStateReq_Level{float32(val)}
            }
        }
        for _, dev := range devices {
            fmt.Println(dev)
            state.IeeeAddr = proto.Uint64(*dev.IeeeAddr);
            req := &remote.Req{
                Type: remote.ReqType_SetDeviceState.Enum(),
                SetDeviceState: state,
            }
            req_data, _ := proto.Marshal(req)
            client.SendBytes(req_data, 0)

            resp_data, _ := client.RecvBytes(0)
            resp := &remote.GetDevicesResp{}
            _ = proto.Unmarshal(resp_data, resp)
            devices := resp.Devices
            fmt.Println(devices)
        }
    }
}
