package main

import (
    "fmt"
    "log"
    zmq "github.com/pebbe/zmq4"
    "github.com/golang/protobuf/proto"
    "github.com/tkln/tradfrid/api"
    "github.com/dyrkin/zigbee-steward"
    "github.com/dyrkin/zigbee-steward/configuration"
    "github.com/dyrkin/zigbee-steward/model"
    "strconv"
)

var devices = map[uint64]*model.Device{}

func addDevice(dev *model.Device) {
    log.Print("Found device: ", dev)
    addr, err := strconv.ParseUint(dev.IEEEAddress, 0, 64)
    if err != nil {
        log.Fatal("ParseUint failed: ", err)
    }
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

    server, err := zmq.NewSocket(zmq.REP)
    if err != nil {
        log.Fatal("NewSocket failed:  ", err)
    }
    defer server.Close()
    server.Bind("tcp://*:5432")

    stewie.Start()
    for {
        req := &remote.Req{}
        log.Print("Waiting for recv req");
        req_data, err := server.RecvBytes(0)
        if err != nil {
            log.Fatal("Recv failed: ", err)
        }
        err = proto.Unmarshal(req_data, req)
        if err != nil {
            log.Fatal("Unmarshal failed: ", err)
        }
        log.Print("Parsed req");

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
            resp_data, err := proto.Marshal(resp);
            if err != nil {
                log.Fatal("Marshal failed: ", err);
            }
            log.Print("Sending resp");
            server.SendBytes(resp_data, 0);
        case remote.ReqType_SetDeviceState:
            resp := &remote.StatusResp{
                Status: proto.Int32(0),
            }
            resp_data, err := proto.Marshal(resp)
            if err != nil {
                log.Fatal("Marshal failed: ", err);
            }
            log.Print("Sending resp");
            server.SendBytes(resp_data, 0);
            go func() {
                addr := fmt.Sprintf("0x%x", *req.SetDeviceState.IeeeAddr)
                local := stewie.Functions().Cluster().Local()
                switch x := req.SetDeviceState.Data.(type) {
                case *remote.SetDeviceStateReq_Onoff:
                    OnOff := local.OnOff()
                    if x.Onoff {
                        OnOff.On(addr, 1)
                    } else {
                        OnOff.Off(addr, 1)
                    }
                case *remote.SetDeviceStateReq_Level:
                    level := uint8(255 * x.Level);
                    local.LevelControl().MoveToLevel(addr, 1, level, 1)
                }
            }()
        }
    }
}
