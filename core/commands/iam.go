package commands

import (
	"encoding/json"
	"fmt"
	"gx/ipfs/QmPTfgFTo9PFr1PvPKyKoeMgBvYPh6cX3aDP7DHKVbnCbi/go-ipfs-cmds"
	"gx/ipfs/QmSP88ryZkHSRn1fnngAaV2Vcn63WUJzAavnRM9CVdU1Ky/go-ipfs-cmdkit"
	"io"

	"encoding/base64"
	"github.com/Harold-the-Axeman/dacc-iam-filesystem/iam"
	"github.com/Harold-the-Axeman/dacc-iam-filesystem/iam/pb"
	"gx/ipfs/QmZFbDTY9jfSBms2MchvYM9oYRbAF19K7Pby47yDBfpPrb/go-cid"
	"gx/ipfs/QmcZSzKEM5yDfpZbeEEZaVmaZ1zXm6JWTbrQZSB8hCVPzk/go-libp2p-peer"
	"gx/ipfs/QmdxUuburamoF6zF9qjeQC4WYcWGbWuRmdLacMEsW8ioD8/gogo-protobuf/proto"
)

// IamCmd ...
var IamCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM CMD",
		ShortDescription: "",
	},
	Subcommands: map[string]*cmds.Command{
		"test":   IamTestCmd,
		"csr":    IamCsrCmd,
		"cot":    IamCotCmd,
		"crt":    IamCrtCmd,
		"cat":    IamCatCmd,
		"get":    IamGetCmd,
		"verify": IamVerifyCmd,
	},
}

//IamCsrCmd
var IamCsrCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM CSR CMD",
		ShortDescription: "",
	},
	Subcommands: map[string]*cmds.Command{
		"create": IamCSRCreateCmd,
		"parse":  IamCSRParseCmd,
	},
}

// IamCSRCreateCmd
var IamCSRCreateCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM CSR Create CMD",
		ShortDescription: "",
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("cid", true, true, "The cid key"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {
		if len(req.Arguments) > 0 && len(req.Arguments[0]) > 0 {
			cidArg := req.Arguments[0]
			n, _ := GetNode(env)
			peerId := n.Identity
			c, err := cid.Decode(cidArg)
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}
			csr := iam.MakeCSR(peerId, c)

			signd, err := n.IAM.MakeMessage(csr)
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}
			out, err := makeOutput(signd)
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}
			cmds.EmitOnce(res, out)
		} else {
			res.SetError(fmt.Errorf("cid required"), cmdkit.ErrNormal)
		}
	},
	Type: IamSignMessageOutput{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeEncoder(jsonEncoder),
	},
}

//IAM COT Parse CMD
var IamCSRParseCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM CSR Parse CMD",
		ShortDescription: "",
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("message", true, true, "The message string"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {
		if len(req.Arguments) > 0 && len(req.Arguments[0]) > 0 {
			message := req.Arguments[0]

			rec := new(iam_pb.MessageCSR)
			ret := make(map[string]interface{})

			_ = proto.Unmarshal(base642bytes(message), rec)
			c, _ := cid.Cast(rec.Cid)
			oid, _ := peer.IDFromBytes(rec.OId)
			ret["cid"] = c.String()
			ret["oid"] = oid.Pretty()
			cmds.EmitOnce(res, ret)

		} else {
			res.SetError(fmt.Errorf("message required"), cmdkit.ErrNormal)
		}
	},
	Type: map[string]interface{}{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeEncoder(jsonEncoder),
	},
}

//IAM COT CMD
var IamCotCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM COT CMD",
		ShortDescription: "",
	},
	Subcommands: map[string]*cmds.Command{
		"create": IamCOTCreateCmd,
		"get":    IamCOTGetCmd,
		"parse":  IamCOTParseCmd,
	},
}

//IamCOTCreateCmd
var IamCOTCreateCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM COT Create CMD",
		ShortDescription: "",
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("token", true, true, "The csr json token").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {
		if len(req.Arguments) > 0 && len(req.Arguments[0]) > 0 {
			jsonToken := req.Arguments[0]

			signed, err := parseJsonToSignedMessage(jsonToken)
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}
			n, _ := GetNode(env)

			peerId, err := peer.IDFromBytes(signed.PeerId)
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}
			vr := n.IAM.VerifyData(signed.Message, signed.Signature, peerId, signed.PublicKey)
			if vr {
				csr := new(iam_pb.MessageCSR)

				_ = proto.Unmarshal(signed.Message, csr)

				cot, err := n.IAM.HandleCSR(csr)
				if err != nil {
					res.SetError(err, cmdkit.ErrNormal)
					return
				}

				signd, err := n.IAM.MakeMessage(cot)
				if err != nil {
					res.SetError(err, cmdkit.ErrNormal)
					return
				}
				c, _ := cid.Cast(csr.Cid)
				n.IAM.SaveCOTLocal(c, signd.Signature)
				out, err := makeOutput(signd)
				if err != nil {
					res.SetError(err, cmdkit.ErrNormal)
					return
				}
				cmds.EmitOnce(res, out)
			} else {
				res.SetError(fmt.Errorf("Verify fail"), cmdkit.ErrNormal)
			}

			//cmds.EmitOnce(res, output)
		} else {
			res.SetError(fmt.Errorf("token required"), cmdkit.ErrNormal)
		}
	},
	Type: IamSignMessageOutput{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeEncoder(jsonEncoder),
	},
}

//IAM COT Get CMD
var IamCOTGetCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM COT Get CMD",
		ShortDescription: "",
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("cid", true, true, "The cid string"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {
		if len(req.Arguments) > 0 && len(req.Arguments[0]) > 0 {
			cidArg := req.Arguments[0]

			c, err := cid.Decode(cidArg)
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}

			n, _ := GetNode(env)
			buf, _ := n.IAM.FindCOTLocal(c)

			ret := make(map[string]interface{})
			ret["cot"] = bytes2base64(buf)
			cmds.EmitOnce(res, ret)

		} else {
			res.SetError(fmt.Errorf("cid required"), cmdkit.ErrNormal)
		}
	},
	Type: map[string]interface{}{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeEncoder(jsonEncoder),
	},
}

//IAM COT Parse CMD
var IamCOTParseCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM COT Parse CMD",
		ShortDescription: "",
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("message", true, true, "The message string"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {
		if len(req.Arguments) > 0 && len(req.Arguments[0]) > 0 {
			cotArg := req.Arguments[0]

			rec := new(iam_pb.MessageCOT)
			ret := make(map[string]interface{})

			_ = proto.Unmarshal(base642bytes(cotArg), rec)
			c, _ := cid.Cast(rec.Cid)
			sid, _ := peer.IDFromBytes(rec.SId)
			oid, _ := peer.IDFromBytes(rec.OId)
			ret["cid"] = c.String()
			ret["sid"] = sid.Pretty()
			ret["oid"] = oid.Pretty()
			cmds.EmitOnce(res, ret)

		} else {
			res.SetError(fmt.Errorf("message required"), cmdkit.ErrNormal)
		}
	},
	Type: map[string]interface{}{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeEncoder(jsonEncoder),
	},
}

//IAM CRT CMD
var IamCrtCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM CRT CMD",
		ShortDescription: "",
	},
	Subcommands: map[string]*cmds.Command{
		"create": IamCRTCreateCmd,
		"parse":  IamCRTParseCmd,
	},
}

// IamCRTCreateCmd
var IamCRTCreateCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM CRT create CMD",
		ShortDescription: "",
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("cid", true, true, "The cid key"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {
		if len(req.Arguments) > 0 && len(req.Arguments[0]) > 0 {
			cidArg := req.Arguments[0]
			n, _ := GetNode(env)
			peerId := n.Identity
			c, err := cid.Decode(cidArg)
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}
			crt := iam.MakeCRT(peerId, peerId, c)
			signd, err := n.IAM.MakeMessage(crt)
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}
			out, err := makeOutput(signd)
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}
			cmds.EmitOnce(res, out)
		} else {
			res.SetError(fmt.Errorf("cid required"), cmdkit.ErrNormal)
		}
	},
	Type: IamSignMessageOutput{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeEncoder(jsonEncoder),
	},
}

// IamCRTParseCmd
var IamCRTParseCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM CRT Parse CMD",
		ShortDescription: "",
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("message", true, true, "The message string"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {
		if len(req.Arguments) > 0 && len(req.Arguments[0]) > 0 {
			message := req.Arguments[0]
			rec := new(iam_pb.MessageCRT)
			ret := make(map[string]interface{})

			_ = proto.Unmarshal(base642bytes(message), rec)
			c, _ := cid.Cast(rec.Cid)
			rid, _ := peer.IDFromBytes(rec.RId)
			oid, _ := peer.IDFromBytes(rec.OId)
			ret["cid"] = c.String()
			ret["rid"] = rid.Pretty()
			ret["oid"] = oid.Pretty()
			cmds.EmitOnce(res, ret)

		} else {
			res.SetError(fmt.Errorf("message required"), cmdkit.ErrNormal)
		}
	},
	Type: map[string]interface{}{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeEncoder(jsonEncoder),
	},
}

//IAM CAT CMD
var IamCatCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM CAT CMD",
		ShortDescription: "",
	},
	Subcommands: map[string]*cmds.Command{
		"create": IamCATCreateCmd,
		"parse":  IamCATParseCmd,
	},
}

// IamCRTCmd
var IamCATCreateCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM CAR Create CMD",
		ShortDescription: "",
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("token", true, true, "The crt json token").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {
		if len(req.Arguments) > 0 && len(req.Arguments[0]) > 0 {

			jsonToken := req.Arguments[0]

			signed, err := parseJsonToSignedMessage(jsonToken)
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}
			n, _ := GetNode(env)

			peerId, err := peer.IDFromBytes(signed.PeerId)
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}
			vr := n.IAM.VerifyData(signed.Message, signed.Signature, peerId, signed.PublicKey)
			if vr {
				crt := new(iam_pb.MessageCRT)

				_ = proto.Unmarshal(signed.Message, crt)

				cot, err := n.IAM.HandleCRT(crt)
				if err != nil {
					res.SetError(err, cmdkit.ErrNormal)
					return
				}

				signd, err := n.IAM.MakeMessage(cot)
				if err != nil {
					res.SetError(err, cmdkit.ErrNormal)
					return
				}
				c, _ := cid.Cast(crt.Cid)
				n.IAM.SaveCATLocal(c, signd.Signature)
				out, err := makeOutput(signd)
				if err != nil {
					res.SetError(err, cmdkit.ErrNormal)
					return
				}
				cmds.EmitOnce(res, out)
			} else {
				res.SetError(fmt.Errorf("Verify fail"), cmdkit.ErrNormal)
			}
		} else {
			res.SetError(fmt.Errorf("token required"), cmdkit.ErrNormal)
		}
	},
	Type: IamSignMessageOutput{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeEncoder(jsonEncoder),
	},
}

var IamCATParseCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM CAT Parse CMD",
		ShortDescription: "",
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("message", true, true, "The message string"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {
		if len(req.Arguments) > 0 && len(req.Arguments[0]) > 0 {
			cotArg := req.Arguments[0]

			rec := new(iam_pb.MessageCAT)
			ret := make(map[string]interface{})

			_ = proto.Unmarshal(base642bytes(cotArg), rec)
			c, _ := cid.Cast(rec.Cid)
			rid, _ := peer.IDFromBytes(rec.RId)
			oid, _ := peer.IDFromBytes(rec.OId)
			ret["cid"] = c.String()
			ret["rid"] = rid.Pretty()
			ret["oid"] = oid.Pretty()
			ret["cot"] = bytes2base64(rec.Cot)
			cmds.EmitOnce(res, ret)

		} else {
			res.SetError(fmt.Errorf("message required"), cmdkit.ErrNormal)
		}
	},
	Type: map[string]interface{}{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeEncoder(jsonEncoder),
	},
}

//IAM VERIFY CMD
var IamVerifyCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM VERIFY CMD",
		ShortDescription: "",
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("message", true, true, "The message string"),
	},
	Options: []cmdkit.Option{
		cmdkit.StringOption("pubkey", "The pubkey used to check the signed message"),
		cmdkit.StringOption("signed", "The signed message"),
		cmdkit.StringOption("peerid", "The peerid for the pubkey"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {
		if len(req.Arguments) > 0 && len(req.Arguments[0]) > 0 {
			cotArg := req.Arguments[0]

			pubkey, _ := req.Options["pubkey"].(string)
			signed, _ := req.Options["signed"].(string)
			peerid, _ := req.Options["peerid"].(string)

			n, _ := GetNode(env)
			pid, err := peer.IDB58Decode(peerid)
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}
			vr := n.IAM.VerifyData(base642bytes(cotArg), base642bytes(signed), pid, base642bytes(pubkey))

			cmds.EmitOnce(res, fmt.Sprintf("Verify result:%t", vr))
		} else {
			res.SetError(fmt.Errorf("cid required"), cmdkit.ErrNormal)
		}
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeEncoder(stringEncoder),
	},
}

//parseJsonToOutputStruct
func parseJsonToOutputStruct(jsonStr string) (*IamSignMessageOutput, error) {
	p := &IamSignMessageOutput{}
	err := json.Unmarshal([]byte(jsonStr), p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

//parseJsonToSignedMessage
func parseJsonToSignedMessage(jsonStr string) (*iam_pb.SignedMessage, error) {
	p := &IamSignMessageOutput{}
	err := json.Unmarshal([]byte(jsonStr), p)
	if err != nil {
		return nil, err
	}
	return parseOutput(p)
}

//parseOutput
func parseOutput(output *IamSignMessageOutput) (*iam_pb.SignedMessage, error) {
	signd := new(iam_pb.SignedMessage)
	pid, err := peer.IDB58Decode(output.PeerId)
	if err != nil {
		return nil, err
	}
	signd.PeerId = []byte(pid)
	signd.Message = base642bytes(output.Message)
	signd.Signature = base642bytes(output.Signature)
	signd.PublicKey = base642bytes(output.PublicKey)
	return signd, nil
}

//make output
func makeOutput(msg *iam_pb.SignedMessage) (*IamSignMessageOutput, error) {
	pid, err := peer.IDFromBytes(msg.PeerId)
	if err != nil {
		return nil, err
	}
	out := IamSignMessageOutput{
		PeerId:    pid.Pretty(),
		Message:   bytes2base64(msg.Message),
		PublicKey: bytes2base64(msg.PublicKey),
		Signature: bytes2base64(msg.Signature),
	}
	return &out, nil
}

//string2bytes
func base642bytes(str string) []byte {
	ret, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Fatal(err)
	}
	return ret
}

//bytes2string
func bytes2base64(bs []byte) string {
	return base64.StdEncoding.EncodeToString(bs)
}

//json encoder
func jsonEncoder(request *cmds.Request, writer io.Writer, i interface{}) error {
	marshaled, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		return err
	}
	marshaled = append(marshaled, byte('\n'))
	_, err = writer.Write(marshaled)
	return err
}

func stringEncoder(request *cmds.Request, writer io.Writer, i interface{}) error {
	s := i.(string)
	ret := []byte(fmt.Sprintf("%s\n", s))
	_, err := writer.Write(ret)
	return err
}

//IamSignMessageOutput
type IamSignMessageOutput struct {
	PeerId    string
	Message   string
	PublicKey string
	Signature string
}

// IAMOutput ...
type IAMOutput struct {
	Content     string
	ContentHash string
	Key         string
}

// ContentOutput ...
type ContentOutput struct {
	Key     string
	Content string
}

// IamGetCmd ...
var IamGetCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM Get CMD",
		ShortDescription: "",
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("cid", true, true, "The cid key"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {
		n, _ := GetNode(env)
		fmt.Println(n.IAM)
		cmds.EmitOnce(res, n.IAM)
	},
}
