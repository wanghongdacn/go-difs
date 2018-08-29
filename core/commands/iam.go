package commands

import (
	"fmt"
	"gx/ipfs/QmPTfgFTo9PFr1PvPKyKoeMgBvYPh6cX3aDP7DHKVbnCbi/go-ipfs-cmds"
	"gx/ipfs/QmPdKqUcHGFdeSpvjVoaTRPPstGif9GBZb5Q56RVw9o69A/go-ipfs-util"
	"gx/ipfs/QmSP88ryZkHSRn1fnngAaV2Vcn63WUJzAavnRM9CVdU1Ky/go-ipfs-cmdkit"
	ds "gx/ipfs/QmVG5gxteQNEMhrS8prJSmU2C9rebtFuTd3SYZ5kE3YZ5k/go-datastore"

	oldcmds "github.com/Harold-the-Axeman/dacc-iam-filesystem/commands"
	lgc "github.com/Harold-the-Axeman/dacc-iam-filesystem/commands/legacy"
	"github.com/Harold-the-Axeman/dacc-iam-filesystem/core"
	b58 "github.com/mr-tron/base58/base58"
)

// IamCmd ...
var IamCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM CMD",
		ShortDescription: "",
	},
	Subcommands: map[string]*cmds.Command{
		"info":    lgc.NewCommand(IamInfoCmd),
		"content": IamContentCmd,
		"test":    IamTestCmd,
		"csr":     lgc.NewCommand(IamInfoCmd),
		"cot":     lgc.NewCommand(IamInfoCmd),
		"crt":     lgc.NewCommand(IamInfoCmd),
		"cat":     lgc.NewCommand(IamInfoCmd),
		"get":     IamGetCmd,
	},
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
