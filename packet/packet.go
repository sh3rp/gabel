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

//
// Represents the function contracts for TLV packets
//

type TLV interface {
	Type() int
	Length() int
	Data() []byte
	ParseFrom([]byte)
	Serialize() []byte
}

//
// Parse and generate an array of TLV's from a Babel packet
//

func ParseBabelPacket(bytes []byte) ([]TLV, error) {
	var tlvs []TLV

	if bytes[0] != 42 {
		return nil, errors.New(fmt.Sprintf("Malformed packet, magic number incorrect (%d)", bytes[0]))
	}

	if bytes[1] != 2 {
		return nil, errors.New(fmt.Sprintf("Packet version unknown (got %d, expected 2)", bytes[1]))
	}

	len := int(bytes[2])<<8 | int(bytes[3])
	endIdx := 4 + len
	currentTLVIdx := 5

	for currentTLVIdx <= endIdx {
		tlvLen := int(bytes[currentTLVIdx+1])<<8 | int(bytes[currentTLVIdx+2])
		var tlv TLV
		switch bytes[currentTLVIdx] {
		case ACKREQ:
			ackReq := new(AckReq)
			ackReq.ParseFrom(bytes[currentTLVIdx : tlvLen+3])
			tlv = ackReq

		case ACK:
			ack := new(Ack)
			ack.ParseFrom(bytes[currentTLVIdx : tlvLen+3])
			tlv = ack

		default:
		}
		tlvs = append(tlvs, tlv)
		currentTLVIdx = currentTLVIdx + 2 + tlvLen
	}

	return tlvs, nil
}

//
// Generate a serialized Babel packet from an array of tlvs
//

func SerializeBabelPacket(tlvlist []TLV) ([]byte, error) {
	var bytes []byte

	totalLen := 0

	if tlvlist != nil {
		for _, tlv := range tlvlist {
			totalLen = totalLen + tlv.Length() + 2
		}
	}

	bytes = make([]byte, 4+totalLen)

	bytes[0] = byte(MAGIC)
	bytes[1] = byte(VERSION)
	bytes[2] = byte(totalLen >> 8)
	bytes[3] = byte(totalLen & 0x00ff)

	currentIdx := 4

	for _, tlv := range tlvlist {
		bytes[currentIdx] = byte(tlv.Type())
		bytes[currentIdx+1] = byte(tlv.Length())
		copy(bytes[currentIdx+1:tlv.Length()], tlv.Data())
	}

	return bytes, nil
}

//
// ACKREQ message codec
//

type AckReq struct {
	Nonce    int16
	Interval int16
}

func (ackReq *AckReq) ParseFrom(bytes []byte) {
	ackReq.Nonce = int16(bytes[4])<<8 | int16(bytes[5])
	ackReq.Interval = int16(bytes[6])<<8 | int16(bytes[7])
}

func (ackReq *AckReq) Serialize() []byte {
	var bytes = make([]byte, 8)

	bytes[0] = byte(ACKREQ)
	bytes[1] = byte(6)
	bytes[2] = 0
	bytes[3] = 0
	bytes[4] = byte(ackReq.Nonce >> 8)
	bytes[5] = byte(ackReq.Nonce & 0x00ff)
	bytes[6] = byte(ackReq.Interval >> 8)
	bytes[7] = byte(ackReq.Interval & 0x00ff)

	return bytes
}

func (ackReq *AckReq) Type() int {
	return ACKREQ
}

func (ackReq *AckReq) Length() int {
	return 6
}

func (ackReq *AckReq) Data() []byte {
	return ackReq.Serialize()[2:]
}

//
// ACK message codec
//

type Ack struct {
	Nonce int16
}

func (ack *Ack) ParseFrom(bytes []byte) {
	ack.Nonce = int16(bytes[2])<<8 | int16(bytes[3])
}

func (ack *Ack) Serialize() []byte {
	var bytes = make([]byte, 4)
	bytes[0] = byte(ACK)
	bytes[1] = 2
	bytes[2] = byte(ack.Nonce >> 8)
	bytes[3] = byte(ack.Nonce & 0x00ff)
	return bytes
}

func (ack *Ack) Type() int {
	return ACK
}

func (ack *Ack) Length() int {
	return 4
}

func (ack *Ack) Data() []byte {
	return ack.Serialize()[2:]
}

//
// HELLO message codec
//

type Hello struct {
	Seqno    int16
	Interval int16
}

func (hello *Hello) ParseFrom(bytes []byte) {}
func (hello *Hello) Serialize() []byte      { return nil }
func (hello *Hello) Type() int              { return 0 }
func (hello *Hello) Length() int            { return 0 }
func (hello *Hello) Data() []byte           { return nil }

//
// IHeardU message codec
//

type IHeardU struct {
	AE      byte
	RxCost  int16
	Address []byte
}

func (ihu *IHeardU) ParseFrom(bytes []byte) {}
func (ihu *IHeardU) Serialize() []byte      { return nil }
func (ihu *IHeardU) Type() int              { return 0 }
func (ihu *IHeardU) Length() int            { return 0 }
func (ihu *IHeardU) Data() []byte           { return nil }

//
// Router-Id message codec
//

type RouterId struct {
	RouterId int64
}

func (routerId *RouterId) ParseFrom(bytes []byte) {}
func (routerId *RouterId) Serialize() []byte      { return nil }
func (routerId *RouterId) Type() int              { return 0 }
func (routerId *RouterId) Length() int            { return 0 }
func (routerId *RouterId) Data() []byte           { return nil }

//
// NextHop message codec
//

type NextHop struct {
	AE      byte
	Address []byte
}

func (nextHop *NextHop) ParseFrom(bytes []byte) {}
func (nextHop *NextHop) Serialize() []byte      { return nil }
func (nextHop *NextHop) Type() int              { return 0 }
func (nextHop *NextHop) Length() int            { return 0 }
func (nextHop *NextHop) Data() []byte           { return nil }

//
// Update message codec
//

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

func (update *Update) ParseFrom(bytes []byte) {}
func (update *Update) Serialize() []byte      { return nil }
func (update *Update) Type() int              { return 0 }
func (update *Update) Length() int            { return 0 }
func (update *Update) Data() []byte           { return nil }

//
// RouteRequest message codec
//

type RouteRequest struct {
	AE        byte
	PrefixLen byte
	Prefix    []byte
}

func (routeRequest *RouteRequest) ParseFrom(bytes []byte) {}
func (routeRequest *RouteRequest) Serialize() []byte      { return nil }
func (routeRequest *RouteRequest) Type() int              { return 0 }
func (routeRequest *RouteRequest) Length() int            { return 0 }
func (routeRequest *RouteRequest) Data() []byte           { return nil }

//
// SeqNo message codec
//

type SeqnoRequest struct {
	AE        byte
	PrefixLen byte
	Seqno     int16
	HopCount  byte
	RouterId  RouterId
	Prefix    []byte
}

func (seqnoRequest *SeqnoRequest) ParseFrom(bytes []byte) {}
func (seqnoRequest *SeqnoRequest) Serialize() []byte      { return nil }
func (seqnoRequest *SeqnoRequest) Type() int              { return 0 }
func (seqnoRequest *SeqnoRequest) Length() int            { return 0 }
func (seqnoRequest *SeqnoRequest) Data() []byte           { return nil }
