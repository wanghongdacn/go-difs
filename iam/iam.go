package iam

import (
	"context"
	host "gx/ipfs/QmRRCrNRs4qxotXx7WJT6SpCvSNEhXvyBcVjXY2K71pcjE/go-libp2p-host"
	ds "gx/ipfs/QmVG5gxteQNEMhrS8prJSmU2C9rebtFuTd3SYZ5kE3YZ5k/go-datastore"
	"log"

	crypto "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"

	inet "gx/ipfs/QmX5J1q63BrrDTbpcHifrFbxH3cMZsvaNajy6u3zCpzBXs/go-libp2p-net"

	//csr "github.com/DecentralizedAccessibleContentChain/dacc-iam-filesystem/iam/pb"
	protocol "gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"

	"github.com/gogo/protobuf/proto"
	"github.com/jbenet/goprocess"
	peer "gx/ipfs/QmcZSzKEM5yDfpZbeEEZaVmaZ1zXm6JWTbrQZSB8hCVPzk/go-libp2p-peer"
)

// ProtocolCSRRequest pattern: /protocol-name/request-or-response-message/version
var ProtocolCSRRequest protocol.ID = "/difs/csrreq/0.0.1"

// ProtocolCSRResponse pattern: /protocol-name/request-or-response-message/version
var ProtocolCSRResponse protocol.ID = "/difs/csrresp/0.0.1"

// IAM ...
type IAM struct {
	host      host.Host
	datastore ds.Datastore // Local data

	ctx  context.Context
	proc goprocess.Process

	protocols []protocol.ID // IAM protocols
}

// NewIAMService ...
func NewIAMService(ctx context.Context, h host.Host, ds ds.Datastore) *IAM {
	iam := &IAM{
		host:      h,
		datastore: ds,
		ctx:       ctx,
	}

	// register for network notifs.
	/*
		iam.host.Network().Notify((*netNotifiee)(iam))

		iam.proc = goprocessctx.WithContextAndTeardown(ctx, func() error {
			// remove ourselves from network notifs.
			iam.host.Network().StopNotify((*netNotifiee)(iam))
			return nil
		})*/

	h.SetStreamHandler(ProtocolCSRRequest, onCSRRequest)   // 以后可以拆出来到 iam_net里面
	h.SetStreamHandler(ProtocolCSRResponse, onCSRResponse) // 以后可以拆出来到 iam_net里面

	return iam
}

// remote peer requests handler
func onCSRRequest(s inet.Stream) {
}

//
func onCSRResponse(s inet.Stream) {

}

// copy from node.go
// Verify incoming p2p message data integrity
// data: data to verify
// signature: author signature provided in the message payload
// peerId: author peer id from the message payload
// pubKeyData: author public key from the message payload
func (i *IAM) VerifyData(data []byte, signature []byte, peerID peer.ID, pubKeyData []byte) bool {
	key, err := crypto.UnmarshalPublicKey(pubKeyData)
	if err != nil {
		log.Println(err, "Failed to extract key from message key data")
		return false
	}

	// extract node id from the provided public key
	idFromKey, err := peer.IDFromPublicKey(key)

	if err != nil {
		log.Println(err, "Failed to extract peer id from public key")
		return false
	}

	// verify that message author node id matches the provided node public key
	if idFromKey != peerID {
		log.Println(err, "Node id and provided public key mismatch")
		return false
	}

	res, err := key.Verify(data, signature)
	if err != nil {
		log.Println(err, "Error authenticating data")
		return false
	}

	return res
}

// copy from node.go , sign an outgoing p2p message payload
func (i *IAM) SignProtoMessage(message proto.Message) ([]byte, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	return i.SignData(data)
}

// copy from node.go. sign binary data using the local node's private key
func (i *IAM) SignData(data []byte) ([]byte, error) {
	key := i.host.Peerstore().PrivKey(i.host.ID())
	res, err := key.Sign(data)
	return res, err
}

/*
// Verify ...
func (msg *Message) Verify() (bool, error) {
	// TODO: 1.check multiHash(PubKey) === ID

	// TODO: 2. include ID
	return msg.PubKey.Verify(msg.Payload, msg.Signature)
}

// Sign ....
func (msg *Message) Sign(privKey crypto.PrivKey) ([]byte, error) {
	//TODO: add ID in the Payload
	return privKey.Sign(msg.Payload)
}
*/
