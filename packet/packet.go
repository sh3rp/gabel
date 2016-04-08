package packet

import (
	"errors"
	"fmt"
)

//TODO: Convert byte array to struct codecs to use encoding/binary?

// magic number for start of a babel packet

var MAGIC = 42

// version for start of a babel packet

var VERSION = 2

// TLV types

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
	var bytes = make([]byte, ackReq.Length()+2)

	bytes[0] = byte(ackReq.Type())
	bytes[1] = byte(ackReq.Length())
	bytes[2] = 0
	bytes[3] = 0
	bytes[4] = byte(ackReq.Nonce >> 8)
	bytes[5] = byte(ackReq.Nonce & 0x00ff)
	bytes[6] = byte(ackReq.Interval >> 8)
	bytes[7] = byte(ackReq.Interval & 0x00ff)

	return bytes
}

func (ackReq *AckReq) Type() int    { return ACKREQ }
func (ackReq *AckReq) Length() int  { return 6 }
func (ackReq *AckReq) Data() []byte { return ackReq.Serialize()[2:] }

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
	var bytes = make([]byte, ack.Length()+2)
	bytes[0] = byte(ack.Type())
	bytes[1] = byte(ack.Length())
	bytes[2] = byte(ack.Nonce >> 8)
	bytes[3] = byte(ack.Nonce & 0x00ff)
	return bytes
}

func (ack *Ack) Type() int    { return ACK }
func (ack *Ack) Length() int  { return 2 }
func (ack *Ack) Data() []byte { return ack.Serialize()[2:] }

//
// HELLO message codec
//

type Hello struct {
	Seqno    int16
	Interval int16
}

func (hello *Hello) ParseFrom(bytes []byte) {
	hello.Seqno = int16(bytes[4])<<8 | int16(bytes[5])
	hello.Interval = int16(bytes[6])<<8 | int16(bytes[7])
}

func (hello *Hello) Serialize() []byte {
	var bytes = make([]byte, hello.Length()+2)
	bytes[0] = byte(hello.Type())
	bytes[1] = byte(hello.Length())
	bytes[2] = 0
	bytes[3] = 0
	bytes[4] = byte(hello.Seqno >> 8)
	bytes[5] = byte(hello.Seqno & 0x00ff)
	bytes[6] = byte(hello.Interval >> 8)
	bytes[7] = byte(hello.Interval & 0x00ff)
	return bytes
}

func (hello *Hello) Type() int    { return HELLO }
func (hello *Hello) Length() int  { return 6 }
func (hello *Hello) Data() []byte { return hello.Serialize()[2:] }

//
// IHeardU message codec
//

type IHeardU struct {
	AE       byte
	RxCost   int16
	Interval int16
	Address  []byte
}

func (ihu *IHeardU) ParseFrom(bytes []byte) {
	ihu.AE = bytes[2]
	// bytes[3] is reserved
	ihu.RxCost = int16(bytes[4])<<8 | int16(bytes[5])
	ihu.Interval = int16(bytes[6])<<8 | int16(bytes[7])
	ihu.Address = bytes[8:]
}

func (ihu *IHeardU) Serialize() []byte {
	var bytes = make([]byte, ihu.Length()+2)

	bytes[0] = byte(ihu.Type())
	bytes[1] = byte(ihu.Length())
	bytes[2] = ihu.AE
	// bytes[3] is reserved
	bytes[4] = byte(ihu.RxCost >> 8)
	bytes[5] = byte(ihu.RxCost & 0x00ff)
	bytes[6] = byte(ihu.Interval >> 8)
	bytes[7] = byte(ihu.Interval & 0x00ff)

	copy(bytes[8:], ihu.Address)

	return bytes
}

func (ihu *IHeardU) Type() int    { return IHU }
func (ihu *IHeardU) Length() int  { return 3 + len(ihu.Address) }
func (ihu *IHeardU) Data() []byte { return ihu.Serialize()[2:] }

//
// Router-Id message codec
//

type RouterId struct {
	RouterId int64
}

func (routerId *RouterId) ParseFrom(bytes []byte) {
	routerId.RouterId = int64(bytes[4])<<56 |
		(int64(bytes[5]) << 48 & 0x00ff000000000000) |
		(int64(bytes[6]) << 40 & 0x0000ff0000000000) |
		(int64(bytes[7]) << 32 & 0x000000ff00000000) |
		(int64(bytes[8]) << 24 & 0x00000000ff000000) |
		(int64(bytes[9]) << 16 & 0x0000000000ff0000) |
		(int64(bytes[10]) << 8 & 0x000000000000ff00) |
		(int64(bytes[11]) & 0x00000000000000ff)
}

func (routerId *RouterId) Serialize() []byte {
	var bytes = make([]byte, routerId.Length()+2)
	bytes[0] = byte(routerId.Type())
	bytes[1] = byte(routerId.Length())
	bytes[2] = 0
	bytes[3] = 0
	bytes[4] = byte(routerId.RouterId >> 56)
	bytes[5] = byte(routerId.RouterId >> 48 & 255)
	bytes[6] = byte(routerId.RouterId >> 40 & 255)
	bytes[7] = byte(routerId.RouterId >> 32 & 255)
	bytes[8] = byte(routerId.RouterId >> 24 & 255)
	bytes[9] = byte(routerId.RouterId >> 16 & 255)
	bytes[10] = byte(routerId.RouterId >> 8 & 255)
	bytes[11] = byte(routerId.RouterId & 255)
	return bytes
}

func (routerId *RouterId) Type() int    { return ROUTERID }
func (routerId *RouterId) Length() int  { return 8 }
func (routerId *RouterId) Data() []byte { return routerId.Serialize()[2:] }

//
// NextHop message codec
// TODO: verify functionality
//

type NextHop struct {
	AE      byte
	Address []byte
}

func (nextHop *NextHop) ParseFrom(bytes []byte) {
	nextHop.AE = bytes[2]
	// bytes[3] is reserved
	nextHop.Address = bytes[4:]
}

func (nextHop *NextHop) Serialize() []byte {
	var bytes = make([]byte, nextHop.Length()+2)

	bytes[0] = byte(nextHop.Type())
	bytes[1] = byte(nextHop.Length())
	bytes[2] = nextHop.AE
	bytes[3] = 0
	copy(bytes[4:], nextHop.Address)

	return bytes
}

func (nextHop *NextHop) Type() int    { return NEXTHOP }
func (nextHop *NextHop) Length() int  { return 2 + len(nextHop.Address) }
func (nextHop *NextHop) Data() []byte { return nextHop.Serialize()[2:] }

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

func (update *Update) ParseFrom(bytes []byte) {
	update.AE = bytes[2]
	update.Flags = bytes[3]
	update.PrefixLen = bytes[4]
	update.Omitted = bytes[5]
	update.Interval = int16(bytes[6])<<8 | int16(bytes[7])
	update.Seqno = int16(bytes[8])<<8 | int16(bytes[9])
	update.Metric = int16(bytes[10])<<8 | int16(bytes[11])
	update.Prefix = make([]byte, update.PrefixLen)
	copy(update.Prefix, bytes[12:])
}

func (update *Update) Serialize() []byte {
	var bytes = make([]byte, update.Length()+2)

	bytes[0] = byte(update.Type())
	bytes[1] = byte(update.Length())
	bytes[2] = update.AE
	bytes[3] = update.Flags
	bytes[4] = update.PrefixLen
	bytes[5] = update.Omitted
	bytes[6] = byte(update.Interval >> 8)
	bytes[7] = byte(update.Interval & 0x00ff)
	bytes[8] = byte(update.Seqno >> 8)
	bytes[9] = byte(update.Seqno & 0x00ff)
	bytes[10] = byte(update.Metric >> 8)
	bytes[11] = byte(update.Metric & 0x00ff)
	copy(bytes[12:], update.Prefix)

	return bytes
}

func (update *Update) Type() int    { return UPDATE }
func (update *Update) Length() int  { return 10 + len(update.Prefix) }
func (update *Update) Data() []byte { return update.Serialize()[2:] }

//
// RouteRequest message codec
//

type RouteRequest struct {
	AE        byte
	PrefixLen byte
	Prefix    []byte
}

func (routeRequest *RouteRequest) ParseFrom(bytes []byte) {
	routeRequest.AE = bytes[2]
	routeRequest.PrefixLen = bytes[3]
	routeRequest.Prefix = make([]byte, routeRequest.PrefixLen)
	copy(routeRequest.Prefix, bytes[4:])
}

func (routeRequest *RouteRequest) Serialize() []byte {
	var bytes = make([]byte, routeRequest.Length()+2)

	bytes[0] = byte(routeRequest.Type())
	bytes[1] = byte(routeRequest.Length())
	bytes[2] = routeRequest.AE
	bytes[3] = routeRequest.PrefixLen
	copy(bytes[4:], routeRequest.Prefix)

	return bytes
}

func (routeRequest *RouteRequest) Type() int    { return ROUTEREQ }
func (routeRequest *RouteRequest) Length() int  { return 2 + len(routeRequest.Prefix) }
func (routeRequest *RouteRequest) Data() []byte { return routeRequest.Serialize()[2:] }

//
// SeqNo message codec
//

type SeqnoRequest struct {
	AE        byte
	PrefixLen byte
	Seqno     int16
	HopCount  byte
	RouterId  int64
	Prefix    []byte
}

func (seqnoRequest *SeqnoRequest) ParseFrom(bytes []byte) {
	seqnoRequest.AE = bytes[3]
	seqnoRequest.PrefixLen = bytes[4]
	seqnoRequest.Seqno = int16(bytes[5])<<8 | int16(bytes[6])&0x00ff
	seqnoRequest.HopCount = bytes[7]
	seqnoRequest.RouterId = int64(bytes[8])<<56 |
		(int64(bytes[9]) << 48 & 0x00ff000000000000) |
		(int64(bytes[10]) << 40 & 0x0000ff0000000000) |
		(int64(bytes[11]) << 32 & 0x000000ff00000000) |
		(int64(bytes[12]) << 24 & 0x00000000ff000000) |
		(int64(bytes[13]) << 16 & 0x0000000000ff0000) |
		(int64(bytes[14]) << 8 & 0x000000000000ff00) |
		(int64(bytes[15]) & 0x00000000000000ff)
	seqnoRequest.Prefix = make([]byte, seqnoRequest.PrefixLen)
	copy(seqnoRequest.Prefix, bytes[16:])
}

func (seqnoRequest *SeqnoRequest) Serialize() []byte {
	var bytes = make([]byte, seqnoRequest.Length()+2)
	return bytes
}

func (seqnoRequest *SeqnoRequest) Type() int    { return SEQNOREQ }
func (seqnoRequest *SeqnoRequest) Length() int  { return 7 + len(seqnoRequest.Prefix) }
func (seqnoRequest *SeqnoRequest) Data() []byte { return seqnoRequest.Serialize()[2:] }
