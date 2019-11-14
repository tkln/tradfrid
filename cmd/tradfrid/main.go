package main

import (
    "log"
    zmq "github.com/pebbe/zmq4"
    "github.com/golang/protobuf/proto"
    "tradfrid/proto"
    "github.com/dyrkin/zigbee-steward"
    "github.com/dyrkin/zigbee-steward/configuration"
    "github.com/dyrkin/zigbee-steward/model"
    "strconv"
)

var devices = map[uint64]*model.Device{}

func addDevice(dev *model.Device) {
    log.Print("Found device: ", dev)
    addr, _ := strconv.ParseUint(dev.IEEEAddress, 0, 64)
    devices[addr] = dev
    log.Print(devices)
}

func main() {
    log.Print("Starting")

    znpConf := configuration.Default()
    znpConf.Serial.PortName = "/dev/ttyACM0"
    znpConf.PermitJoin = true

    stewie := steward.New(znpConf)

    handleZnpEvent := func() {
        for {
            select {
            case dev:= <-stewie.Channels().OnDeviceBecameAvailable():
                addDevice(dev)
            case dev:= <-stewie.Channels().OnDeviceRegistered():
                addDevice(dev)
            }
        }
    }

    go handleZnpEvent()

    server, _ := zmq.NewSocket(zmq.REP)
    defer server.Close()
    server.Bind("tcp://*:5432")

    stewie.Start()
    for {
        log.Print("Req");
        req := &remote.Req{}
        req_data, _ := server.RecvBytes(0)
        _ = proto.Unmarshal(req_data, req)

        switch *req.Type {
        case remote.ReqType_GetDevices:
            log.Print("Get devices req");
            resp := &remote.GetDevicesResp{}
            for k, v := range devices {
                dev := &remote.Device{
                    IeeeAddr: proto.Uint64(k),
                    Model: proto.String(v.Model),
                }
                log.Print(dev)
                resp.Devices = append(resp.Devices, dev)
            }
            resp_data, _ := proto.Marshal(resp);
            server.SendBytes(resp_data, 0);
        }
    }
}
