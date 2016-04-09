package packet

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/skendall/gabel/packet"
)

func TestAckReq(t *testing.T) {
	ackReq := new(packet.AckReq)
	ackReq.Nonce = 1024
	ackReq.Interval = 2048

	bytes := ackReq.Serialize()
	fmt.Println(bytes)
	newAckReq := new(packet.AckReq)
	newAckReq.ParseFrom(bytes)

	assert.Equal(t, int16(1024), newAckReq.Nonce)
	assert.Equal(t, int16(2048), newAckReq.Interval)
}

func TestAck(t *testing.T) {
	ack := new(packet.Ack)
	ack.Nonce = 1024

	bytes := ack.Serialize()
	newAck := new(packet.Ack)
	newAck.ParseFrom(bytes)

	assert.Equal(t, int16(1024), ack.Nonce)
}
