package iam

import (
	"github.com/DecentralizedAccessibleContentChain/dacc-iam-filesystem/iam/pb"

	"fmt"
	"github.com/gogo/protobuf/proto"
	"gx/ipfs/QmQsErDt8Qgw1XrsXf2BpEzDgGWtB1YLsTAARBup5b6B9W/go-libp2p-peer"
	ds "gx/ipfs/QmVG5gxteQNEMhrS8prJSmU2C9rebtFuTd3SYZ5kE3YZ5k/go-datastore"
	"gx/ipfs/QmZFbDTY9jfSBms2MchvYM9oYRbAF19K7Pby47yDBfpPrb/go-cid"
)

//MakeCSR
func MakeCSR(oid peer.ID, cid *cid.Cid) *iam_pb.MessageCSR {
	msg := new(iam_pb.MessageCSR)
	msg.Cid = []byte(cid.Bytes())
	msg.OId = []byte(oid)
	return msg
}

// HandleCSR return COT
func (iam *IAM) HandleCSR(csr *iam_pb.MessageCSR) (*iam_pb.MessageCOT, error) {
	msg := new(iam_pb.MessageCOT)
	msg.OId = csr.OId
	msg.Cid = csr.Cid
	msg.SId = []byte(iam.host.ID())
	return msg, nil
}

//MakeCRT
func MakeCRT(rid peer.ID, oid peer.ID, cid *cid.Cid) *iam_pb.MessageCRT {
	msg := new(iam_pb.MessageCRT)
	msg.OId = []byte(oid)
	msg.Cid = []byte(cid.Bytes())
	msg.RId = []byte(rid)
	return msg
}

// HandleCRT,return CAT
func (iam *IAM) HandleCRT(crt *iam_pb.MessageCRT) (*iam_pb.MessageCAT, error) {

	c, err := cid.Cast(crt.Cid)
	if err != nil {
		return nil, err
	}
	bs, err := iam.FindCOTLocal(c)

	if err != nil {
		return nil, err
	}

	msg := new(iam_pb.MessageCAT)
	msg.OId = []byte(iam.host.ID())
	msg.Cid = crt.Cid
	msg.RId = crt.RId
	msg.Cot = bs
	return msg, nil
}

// MakeRT
func MakeRT(rid string, cid *cid.Cid, cat []byte) *iam_pb.MessageRT {
	msg := new(iam_pb.MessageRT)
	msg.Cat = cat
	msg.Cid = []byte(cid.Bytes())
	msg.RId = []byte(rid)
	return msg
}

// HandleRT,return content
func (iam *IAM) HandleRT(rt *iam_pb.MessageRT) {

}

//SaveCOTLocal
func (iam *IAM) SaveCOTLocal(cid *cid.Cid, cotBytes []byte) {
	kf := fmt.Sprintf("/iam/cot/%s", cid.String())
	key := ds.NewKey(kf)
	iam.datastore.Put(key, cotBytes)
}

//FindCOTLocal
func (iam *IAM) FindCOTLocal(cid *cid.Cid) ([]byte, error) {
	kf := fmt.Sprintf("/iam/cot/%s", cid.String())
	key := ds.NewKey(kf)
	return iam.datastore.Get(key)
}

//SaveCATLocal
func (iam *IAM) SaveCATLocal(cid *cid.Cid, cotBytes []byte) {
	kf := fmt.Sprintf("/iam/cat/%s", cid.String())
	key := ds.NewKey(kf)
	iam.datastore.Put(key, cotBytes)
}

//FindCATLocal
func (iam *IAM) FindCATLocal(cid *cid.Cid) ([]byte, error) {
	kf := fmt.Sprintf("/iam/cat/%s", cid.String())
	key := ds.NewKey(kf)
	return iam.datastore.Get(key)
}

// MakeMessage
func (iam *IAM) MakeMessage(message proto.Message) (*iam_pb.SignedMessage, error) {
	signd := new(iam_pb.SignedMessage)
	signd.PeerId = []byte(iam.host.ID())
	data, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	signd.Message = data
	bs, _ := iam.host.Peerstore().PubKey(iam.host.ID()).Bytes()
	signd.PublicKey = bs
	signature, _ := iam.SignData(data)
	signd.Signature = signature

	return signd, nil
}
