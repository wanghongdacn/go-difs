package iam

import peer "github.com/libp2p/go-libp2p-peer"

type Message struct {
	ID      peer.ID
	Payload string
}
