package packet

type Transport interface {
	Send([]byte) error
}
