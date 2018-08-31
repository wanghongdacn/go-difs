package commands

import (
	"fmt"
	oldcmds "github.com/DecentralizedAccessibleContentChain/dacc-iam-filesystem/commands"
	lgc "github.com/DecentralizedAccessibleContentChain/dacc-iam-filesystem/commands/legacy"
	"github.com/DecentralizedAccessibleContentChain/dacc-iam-filesystem/core"
	"github.com/DecentralizedAccessibleContentChain/dacc-iam-filesystem/iam/pb"
	b58 "github.com/mr-tron/base58/base58"
	"gx/ipfs/QmPTfgFTo9PFr1PvPKyKoeMgBvYPh6cX3aDP7DHKVbnCbi/go-ipfs-cmds"
	"gx/ipfs/QmPdKqUcHGFdeSpvjVoaTRPPstGif9GBZb5Q56RVw9o69A/go-ipfs-util"
	ft "gx/ipfs/QmQjEpRiwVvtowhq69dAtB4jhioPVFXiCcWZm9Sfgn7eqc/go-unixfs"
	"gx/ipfs/QmQsErDt8Qgw1XrsXf2BpEzDgGWtB1YLsTAARBup5b6B9W/go-libp2p-peer"
	dag "gx/ipfs/QmRiQCJZ91B7VNmLvA6sxzDuBJGSojS3uXHHVuNr3iueNZ/go-merkledag"
	"gx/ipfs/QmSP88ryZkHSRn1fnngAaV2Vcn63WUJzAavnRM9CVdU1Ky/go-ipfs-cmdkit"
	ds "gx/ipfs/QmVG5gxteQNEMhrS8prJSmU2C9rebtFuTd3SYZ5kE3YZ5k/go-datastore"
	ipld "gx/ipfs/QmX5CsuHyVZeTLxgRSYkgLSDQKb9UjE8xnhQzCEJWWWFsC/go-ipld-format"
	proto "gx/ipfs/QmdxUuburamoF6zF9qjeQC4WYcWGbWuRmdLacMEsW8ioD8/gogo-protobuf/proto"
	base32 "gx/ipfs/QmfVj3x4D6Jkq9SEoi5n2NmoUomLwoeiwnYz2KQa15wRw6/base32"
)

// IamInfoCmd ...
var IamInfoCmd = &oldcmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM CMD",
		ShortDescription: "iam quick start",
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("info", true, true, "The info key"),
	},
	Options: []cmdkit.Option{
		cmdkit.StringOption("key", "TODO"),
	},
	Run: func(req oldcmds.Request, res oldcmds.Response) {

		fmt.Printf("load info ...\n")
		//nd, err := req.InvocContext().GetNode()
		//if err != nil {
		//	res.SetError(err, cmdkit.ErrNormal)
		//	return
		//}
		iam := &IAMOutput{}
		if len(req.Arguments()) > 0 {
			content := req.Arguments()[0]
			iam.Content = content
			iam.ContentHash = b58.Encode(util.Hash([]byte(content)))
		}

		iam.Key, _, _ = req.Option("key").String()
		res.SetOutput(iam.ContentHash)
	},
}

var IamContentCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM content CMD",
		ShortDescription: "",
	},
	Subcommands: map[string]*cmds.Command{
		"put": IamContentPutCmd,
		"get": IamContentGetCmd,
	},
}

var IamContentPutCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM put content local",
		ShortDescription: "",
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("content", true, true, "The content to be added."),
	},
	Options: []cmdkit.Option{
		cmdkit.StringOption("key", "The key of content,default _"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {

		n, _ := GetNode(env)

		fmt.Printf("%s\n", n.Identity)

		iam := &IAMOutput{}

		iam.Key, _ = req.Options["key"].(string)

		if len(req.Arguments) > 0 && len(req.Arguments[0]) > 0 {
			content := req.Arguments[0]
			iam.Content = content
			iam.ContentHash = b58.Encode(util.Hash([]byte(content)))
			key := iam.Key
			if len(key) <= 0 {
				key = "_"
				iam.Key = key
			}
			putIamDatastore(n, key, []byte(content))
		}

		cmds.EmitOnce(res, &iam)
	},
	Type: IAMOutput{},
}

// IamContentGetCmd
var IamContentGetCmd = &cmds.Command{

	Helptext: cmdkit.HelpText{
		Tagline:          "IAM get local content",
		ShortDescription: "",
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("key", true, true, "The key of content."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {

		n, _ := GetNode(env)

		fmt.Printf("%s\n", n.Identity)
		fmt.Println("%s", n.IAM)
		fmt.Printf("%s\n", n.Identity)

		v := &ContentOutput{}
		v.Key = ""
		if len(req.Arguments) > 0 && len(req.Arguments[0]) > 0 {
			key := req.Arguments[0]
			v.Key = key
			v.Content = getIamDatastore(n, key)
		}
		cmds.EmitOnce(res, &v)
	},
	Type: ContentOutput{},
}

/**
	pug value and key into iamdatastore
	add iam datastore first
	1,add new datastore config into /path/to/repo/config
		{
          "child": {
            "compression": "none",
            "path": "iamdatastore",
            "type": "levelds"
          },
          "mountpoint": "/iam",
          "prefix": "iam.leveldb.datastore",
          "type": "measure"
        }
	2,change file /path/to/repo/datastore_spec
	{"mounts":[{"mountpoint":"/iam","path":"iamdatastore","type":"levelds"},{"mountpoint":"/blocks","path":"blocks","shardFunc":"/repo/flatfs/shard/v1/next-to-last/2","type":"flatfs"},{"mountpoint":"/","path":"datastore","type":"levelds"}],"type":"mount"}
*/
func putIamDatastore(node *core.IpfsNode, key string, value []byte) {
	k := fmt.Sprintf("/iam/%s", key)
	gk := ds.NewKey(k)
	node.Repo.Datastore().Put(gk, []byte(value))
}

//get value by key from iamdatastore
func getIamDatastore(node *core.IpfsNode, key string) string {
	k := fmt.Sprintf("/iam/%s", key)
	gk := ds.NewKey(k)
	v, _ := node.Repo.Datastore().Get(gk)
	s := string(v)
	return s
}

//// GetNode extracts the node from the environment.
//func GetNode(env interface{}) (*core.IpfsNode, error) {
//	ctx, ok := env.(*commands.Context)
//	if !ok {
//		return nil, fmt.Errorf("expected env to be of type %T, got %T", ctx, env)
//	}
//	return ctx.GetNode()
//}

var IamTestCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM test cmd",
		ShortDescription: "",
	},
	Subcommands: map[string]*cmds.Command{
		"tt":      IamAutoRunCmd,
		"cid":     IamTestCidCmd,
		"content": IamContentCmd,
		"hash":    lgc.NewCommand(IamInfoCmd),
	},
}

var IamAutoRunCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "",
		ShortDescription: "",
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {

		n, _ := GetNode(env)
		fmt.Printf("iam auto run test %s\n", n.Identity.Pretty())

		id, _ := peer.IDB58Decode("QmQ4PSDphDBpFcRRhK9E87Qi2PyJGuXGeuLLnvpBt5Tavg")
		k := "/pk/" + string(id)
		key := ds.NewKey(base32.RawStdEncoding.EncodeToString([]byte(k)))

		cid := "xxx"
		msg := new(iam_pb.MessageCSR)
		msg.Cid = []byte(cid)

		data, _ := proto.Marshal(msg)
		n.Repo.Datastore().Put(key, data)

		buf, _ := n.Repo.Datastore().Get(key)
		rec := new(iam_pb.MessageCSR)
		_ = proto.Unmarshal(buf, rec)
		fmt.Printf("pb data cid %s\n", string(rec.Cid))

	},
}

var IamTestCidCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM test cid",
		ShortDescription: "",
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("filecontent", true, true, "The content to be added."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {

		//n, _ := GetNode(env)
		if len(req.Arguments) > 0 && len(req.Arguments[0]) > 0 {
			filecontent := req.Arguments[0]
			GetIpldNode([]byte(filecontent))
		} else {
			fmt.Println("len(arg) must gt zero")
		}
	},
}

func GetIpldNode(fileData []byte) (ipld.Node, error) {

	pn := new(dag.ProtoNode)
	pn.SetCidBuilder(dag.V0CidPrefix())

	fsn := ft.NewFSNode(ft.TFile)
	fsn.SetData(fileData)

	fbyte, err := fsn.GetBytes()
	if err != nil {
		return nil, err
	}
	pn.SetData(fbyte)

	pn.EncodeProtobuf(false)

	cid := pn.Cid()
	fmt.Printf("File:%s\n", string(fileData))
	fmt.Printf("GenCid:%s\n", cid.String())

	return nil, err
}

//k := fmt.Sprintf("/iam/%s", string(p.ID))
//gk := ds.NewKey(k)
//v := "["
//for _, md := range p.Addrs {
//d := md.String()
//v = fmt.Sprintf("%s%s,", v, d)
//}
//v += "]"
//n.Repo.Datastore().Put(gk, []byte(v))
