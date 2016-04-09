package core

import "github.com/sh3rp/gabel/packet"

type Transport interface {
	Send(*packet.BabelPacket)
	AddListener(*BabelPacketListener)
}

type BabelPacketListener interface {
	Received(*packet.BabelPacket)
}

type Loopback struct {
	listeners []BabelPacketListener
}

func (loopback *Loopback) Send(packet *packet.BabelPacket) {
	for _, lo := range loopback.listeners {
		lo.Received(packet)
	}
}

func (loopback *Loopback) AddListener(listener BabelPacketListener) {
	loopback.listeners = append(loopback.listeners, listener)
}
