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

	packet := packet.NewBabelPacket().AddTLV(new(packet.Hello)).AddTLV(new(packet.AckReq)).AddTLV(new(packet.RouterId))

	t1.Send(packet)

}
