package core

import (
	"math/rand"
	"sync"
	"time"

	"github.com/sh3rp/gabel/packet"
)

type Interface struct {
	Label        string
	HelloSeqNo   int16
	transport    Transport
	pendingTLVs  chan packet.TLV
	queueLock    *sync.Mutex
	queuePackets bool
}

func NewInterface(label string, transport *Transport) *Interface {
	return &Interface{
		Label:        label,
		HelloSeqNo:   0,
		pendingTLVs:  make(chan packet.TLV, 10),
		queueLock:    &sync.Mutex{},
		queuePackets: true,
	}
}

func (intf *Interface) Send(tlv packet.TLV) error {
	intf.queueLock.Lock()
	intf.pendingTLVs <- tlv
	intf.queueLock.Unlock()
	return nil
}

func (intf *Interface) Start() {
	go func() {
		for intf.queuePackets {
			intf.queueLock.Lock()
			if len(intf.pendingTLVs) > 0 {
				p := packet.NewBabelPacket()
				for t := range intf.pendingTLVs {
					p.AddTLV(t)
				}
				intf.transport.Send(p)
			}
			intf.queueLock.Unlock()
			time.Sleep(intf.jitter())
		}
	}()
}

func (intf *Interface) Stop() {
	intf.queuePackets = false
}

func (intf *Interface) jitter() time.Duration {
	return time.Duration(rand.Intn(1000))
}
