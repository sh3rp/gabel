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

func (state *State) Received(packet *packet.BabelPacket) {
	log.Println("[PACKET] ", packet)
}
