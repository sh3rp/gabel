package core

import "sync"

type Interface struct {
	Label       string
	HelloSeqNo  int16
	transport   *Transport
	pendingTLVs []TLV
	queueLock   *sync.Mutex
}

func NewInterface(label string, transport *Transport) *Interface {
	return &Interface{
		Label:      label,
		HelloSeqNo: 0,
		queueLock:  &sync.Mutex{},
	}
}

func (intf *Interface) Send(tlv TLV) error {
	intf.queueLock.Lock()
	pendingTLVs = append(pendingTLVs, tlv)
	intf.queueLock.Unlock()
	return nil
}

func (intf *Interface) Start() {

}
