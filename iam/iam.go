package iam

import (
	ds "gx/ipfs/QmVG5gxteQNEMhrS8prJSmU2C9rebtFuTd3SYZ5kE3YZ5k/go-datastore"

	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-peer"
)

// Message ..
type Message struct {
	ID        peer.ID
	Payload   []byte
	Signature []byte
	PubKey    crypto.PubKey
}

// IAM ...
type IAM struct {
	datastore ds.Datastore // Local data
}

func (msg *Message) Verify() (bool, error) {
	// TODO: 1.check multiHash(PubKey) === ID

	// TODO: 2. include ID
	return msg.PubKey.Verify(msg.Payload, msg.Signature)
}

func (msg *Message) Sign(privKey crypto.PrivKey) ([]byte, error) {
	//TODO: add ID in the Payload
	return privKey.Sign(msg.Payload)
}
