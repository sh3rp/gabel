package core

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/sh3rp/gabel/packet"
)

type Interface struct {
	Label        string
	HelloSeqNo   int16
	transport    Transport
	pendingTLVs  []interface{}
	queueLock    *sync.Mutex
	queuePackets bool
}

func NewInterface(label string, transport Transport) *Interface {
	return &Interface{
		Label:        label,
		HelloSeqNo:   0,
		transport:    transport,
		queueLock:    &sync.Mutex{},
		queuePackets: true,
	}
}

func (intf *Interface) Received(bytes []byte) {
	p := packet.BabelPacket{}
	err := p.ParseFrom(bytes)

	if err != nil {
		log.Println("ERROR: %v")
	}

	log.Printf("[%s] [PACKET RECV]", intf.Label)
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

func (intf *Interface) Send(tlv interface{}) error {
	intf.queueLock.Lock()
	intf.pendingTLVs = append(intf.pendingTLVs, tlv)
	intf.queueLock.Unlock()
	return nil
}

func (intf *Interface) Start() {
	go func(i *Interface) {
		for i.queuePackets == true {
			i.queueLock.Lock()

			if i.pendingTLVs != nil || len(i.pendingTLVs) > 0 {
				p := packet.NewBabelPacket()
				for _, t := range intf.pendingTLVs {
					p.AddTLV(t)
				}
				if len(p.TLVs) > 0 {
					intf.transport.Send(p.Serialize())
				}
				i.pendingTLVs = nil
			}
			i.queueLock.Unlock()

			time.Sleep(intf.jitter())
		}
	}(intf)
}

func (intf *Interface) Stop() {
	intf.queuePackets = false
}

func (intf *Interface) jitter() time.Duration {
	return time.Duration(rand.Intn(1000)) * time.Millisecond
}
