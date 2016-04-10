package core

import (
	"log"

	"github.com/sh3rp/gabel/packet"
)

type State struct {
	NodeSeqno  int64
	Interfaces []*Interface
	Neighbors  []*Neighbor
	Sources    map[int]*SourceInfo
	Routes     map[int]*Route
	Pending    []*Request
}

func NewState() *State {
	return &State{
		NodeSeqno:  0,
		Interfaces: make([]*Interface, 1),
		Neighbors:  make([]*Neighbor, 1),
		Sources:    make(map[int]*SourceInfo),
		Routes:     make(map[int]*Route),
		Pending:    make([]*Request, 1),
	}
}

func (state *State) NewInterface(ifLabel string) {

}

func (state *State) Received(bytes []byte) {
	p := packet.BabelPacket{}
	err := p.ParseFrom(bytes)

	if err != nil {
		log.Println("ERROR: %v")
	}

	log.Println("[PACKET RECV]")
	for _, t := range p.TLVs {
		switch tlv := t.(type) {
		case *packet.Hello:
			log.Printf("  HELLO: Sequence: %d, Interval: %d", tlv.Seqno, tlv.Interval)
		case *packet.AckReq:
			log.Printf("  ACK-REQUEST: Nonce: %d, Interval: %d", tlv.Nonce, tlv.Interval)
		case *packet.RouterId:
			log.Printf("  Router-ID: %d", tlv.RouterId)
		default:
			log.Println("  UNKNOWN TYPE", tlv)
		}
	}
}
