package core

import (
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
				var pending []packet.TLV
				for t := range intf.pendingTLVs {
					pending = append(pending, t)
				}
				p := packet.NewBabelPacket(pending)
				intf.transport.Send(p)
			}
			intf.queueLock.Unlock()
			time.Sleep(100)
		}
	}()
}

func (intf *Interface) Stop() {
	intf.queuePackets = false
}
