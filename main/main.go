package main

import (
	"github.com/sh3rp/gabel/core"
	"github.com/sh3rp/gabel/packet"
)

func main() {
	t1 := core.Loopback{}
	t2 := core.Loopback{}
	t1state := core.NewState()
	t2state := core.NewState()

	t1.AddListener(t1state)
	t2.AddListener(t2state)

	hello := new(packet.Hello)
	hello.Interval = 15
	hello.Seqno = 31

	ackreq := new(packet.AckReq)
	ackreq.Nonce = 832
	ackreq.Interval = 10

	routerId := new(packet.RouterId)
	routerId.RouterId = 2147483647

	packet := packet.NewBabelPacket().AddTLV(hello).AddTLV(ackreq).AddTLV(routerId)

	t1.Send(packet.Serialize())

}
