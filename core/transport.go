package core

type Transport interface {
	Send(*Packet)
	AddListener(*PacketListener)
}

type PacketListener interface {
	Received(*Packet)
}

type Loopback struct {
	listeners []*PacketListener
}

func (loopback *Loopback) Send(packet *Packet) {
	for _, lo := range loopback.Listeners {
		lo.Received(packet)
	}
}

func (loopback *Loopback) AddListener(listener *PacketListener) {
	loopback.listeners = append(lookback.Listeners, listener)
}
