package core

type Transport interface {
	Send([]byte)
	AddListener(*BabelPacketListener)
}

type BabelPacketListener interface {
	Received([]byte)
}

type Loopback struct {
	listeners []BabelPacketListener
}

func (loopback *Loopback) Send(bytes []byte) {
	for _, lo := range loopback.listeners {
		lo.Received(bytes)
	}
}

func (loopback *Loopback) AddListener(listener BabelPacketListener) {
	loopback.listeners = append(loopback.listeners, listener)
}
