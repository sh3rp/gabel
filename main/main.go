package main

import (
	"time"

	"github.com/sh3rp/gabel/core"
	"github.com/sh3rp/gabel/packet"
)

func main() {
	lo1 := &core.Loopback{}
	lo2 := &core.Loopback{}

	intf1 := core.NewInterface("intf1", lo1)
	intf1.Start()
	intf2 := core.NewInterface("intf2", lo2)
	intf2.Start()

	lo1.AddListener(intf2)
	lo2.AddListener(intf1)

	hello := new(packet.Hello)
	hello.Interval = 15
	hello.Seqno = 31

	intf1.Send(hello)

	ackreq := new(packet.AckReq)
	ackreq.Nonce = 832
	ackreq.Interval = 10

	intf1.Send(ackreq)

	routerId := new(packet.RouterId)
	routerId.RouterId = 2147483647

	intf1.Send(routerId)

	time.Sleep(time.Duration(5000) * time.Millisecond)

}
