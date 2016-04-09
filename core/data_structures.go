package core

import "git.soma.salesforce.com/skendall/gabel/packet"

type State struct {
	NodeSeqno  int64
	Interfaces []*Interface
	Neighbors  []*Neighbor
	Sources    map[int]*SourceInfo
	Routes     map[int]*Route
	Pending    []*Request
}

type Interface struct {
	Label      string
	HelloSeqNo int16
}

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
