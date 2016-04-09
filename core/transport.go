package core

type Transport interface {
	Send(*Packet)
	AddListener(*PacketListener)
}

type PacketListener interface {
	Received(*Packet)
}
