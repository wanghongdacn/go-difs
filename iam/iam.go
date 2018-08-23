package iam

import (
	ds "gx/ipfs/QmVG5gxteQNEMhrS8prJSmU2C9rebtFuTd3SYZ5kE3YZ5k/go-datastore"

	peer "github.com/libp2p/go-libp2p-peer"
)

// Message ..
type Message struct {
	ID      peer.ID
	Payload string
}

// IAM ...
type IAM struct {
	datastore ds.Datastore // Local data
}
