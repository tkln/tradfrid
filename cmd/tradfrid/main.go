package main

import (
    "log"
    zmq "github.com/pebbe/zmq4"
    "github.com/golang/protobuf/proto"
    "tradfrid/proto"
)

func main() {
    log.Print("Starting")
    server, _ := zmq.NewSocket(zmq.REP)
    defer server.Close()
    server.Bind("tcp://*:5432")

    for {
        req := &arith.Req{}
        req_data, _ := server.RecvBytes(0)
        _ = proto.Unmarshal(req_data, req)
        log.Print("Recv: ", req)

        resp := &arith.Resp{}
        if *req.Op == arith.OP_ADD {
            resp.R = proto.Uint32(*req.A + *req.B)
        } else {
            resp.R = proto.Uint32(*req.A - *req.B)
        }

        resp_data, _ := proto.Marshal(resp);
        server.SendBytes(resp_data, 0)
        log.Print("Sent: ", resp)
    }
}
