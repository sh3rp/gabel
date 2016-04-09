package packet

type Transport interface {
	Send(packet *BabelPacket) error
}
