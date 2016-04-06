package packet

import (
	"errors"
	"fmt"
)

var MAGIC = 42
var VERSION = 2

const (
	PAD1     = 0
	PADN     = 1
	ACKREQ   = 2
	ACK      = 3
	HELLO    = 4
	IHU      = 5
	ROUTERID = 6
	NEXTHOP  = 7
	UPDATE   = 8
	ROUTEREQ = 9
	SEQNOREQ = 10
)

type Packet struct {
	Type   int
	Length int
	Data   []byte
}

type AckReq struct {
	Nonce    int16
	Interval int16
}

type Ack struct {
	Nonce int16
}

type Hello struct {
	Seqno    int16
	Interval int16
}

type IHeardU struct {
	AE      byte
	RxCost  int16
	Address []byte
}

type RouterId []byte

type NextHop struct {
	AE      byte
	Address []byte
}

type Update struct {
	AE        byte
	Flags     byte
	PrefixLen byte
	Omitted   byte
	Interval  int16
	Seqno     int16
	Metric    int16
	Prefix    []byte
}

type RouteRequest struct {
	AE        byte
	PrefixLen byte
	Prefix    []byte
}

type SeqnoRequest struct {
	AE        byte
	PrefixLen byte
	Seqno     int16
	HopCount  byte
	RouterId  RouterId
	Prefix    []byte
}

func ParsePacket(bytes []byte) (*Packet, error) {

	if bytes[0] != 42 {
		return nil, errors.New(fmt.Sprintf("Malformed packet, magic number incorrect (%d)", bytes[0]))
	}

	if bytes[1] != 2 {
		return nil, errors.New(fmt.Sprintf("Packet version unknown (got %d, expected 2)", bytes[1]))
	}

	packet := Packet{
		Type:   int(bytes[4]),
		Length: int(bytes[5]),
		Data:   bytes[6:],
	}

	return &packet, nil
}

func ParseAckReq(bytes []byte) (*AckReq, error) {
	var ackReq *AckReq

	ackReq = &AckReq{
		Nonce:    int16(bytes[2]<<8 | bytes[3]),
		Interval: int16(bytes[4]<<8 | bytes[5]),
	}

	return ackReq, nil
}

func ParseAck(bytes []byte) (*Ack, error) {
	var ack *Ack

	ack = &Ack{
		Nonce: int16(bytes[2]<<8 | bytes[3]),
	}

	return ack, nil
}
