package core

import "github.com/sh3rp/gabel/packet"

type Neighbor struct {
	Interface        *Interface
	Address          []byte
	HelloHistory     map[int64]packet.Hello
	TransmissionCost int16
	HelloSeqno       int16
}

type SourceInfo struct {
	Source *Source
	Seqno  int16
	Metric int16
}

type Route struct {
	Source   *Source
	Neighbor *Neighbor
	Metric   int16
	Seqno    int16
	NextHop  packet.NextHop
	Selected bool
}

type Request struct {
}

type Source struct {
	Prefix    []byte
	PrefixLen byte
	RouterId  int64
}

//
// hashing algorithm used to create the index for the map
// based on Java's String.hashCode() algorithm
//

func (source *Source) Hash() int {
	var hash int

	if len(source.Prefix) > 0 {
		for _, b := range source.Prefix {
			hash = int(b*31) + hash
		}
	}

	hash = int(source.PrefixLen*31) + hash
	hash = int(source.RouterId*31) + hash

	return hash
}
