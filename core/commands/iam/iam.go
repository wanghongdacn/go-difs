package iam

import (
	"fmt"
	oldcmds "github.com/Harold-the-Axeman/dacc-iam-filesystem/commands"
	lgc "github.com/Harold-the-Axeman/dacc-iam-filesystem/commands/legacy"
	b58 "github.com/mr-tron/base58/base58"
	"gx/ipfs/QmPTfgFTo9PFr1PvPKyKoeMgBvYPh6cX3aDP7DHKVbnCbi/go-ipfs-cmds"
	"gx/ipfs/QmPdKqUcHGFdeSpvjVoaTRPPstGif9GBZb5Q56RVw9o69A/go-ipfs-util"
	"gx/ipfs/QmSP88ryZkHSRn1fnngAaV2Vcn63WUJzAavnRM9CVdU1Ky/go-ipfs-cmdkit"
)

type IAMOutput struct {
	Key     string
	KeyHash string
	Ih      string
	Iv      string
}

var IamCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM CMD",
		ShortDescription: "",
	},
	Subcommands: map[string]*cmds.Command{
		"info": lgc.NewCommand(IamInfoCmd),
		"new":  IamInfoNewCmd,
		"csr":  lgc.NewCommand(IamInfoCmd),
		"cot":  lgc.NewCommand(IamInfoCmd),
		"crt":  lgc.NewCommand(IamInfoCmd),
		"cat":  lgc.NewCommand(IamInfoCmd),
		"get":  lgc.NewCommand(IamInfoCmd),
	},
}

var IamInfoCmd = &oldcmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM CMD",
		ShortDescription: "iam quick start",
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("key", true, true, "The path to the IPFS object(s) to be outputted."),
	},
	Options: []cmdkit.Option{
		cmdkit.StringOption("ih", "Byte offset to begin reading from."),
		cmdkit.StringOption("iv", "Maximum number of bytes to read."),
	},
	Run: func(req oldcmds.Request, res oldcmds.Response) {

		fmt.Printf("load info ...\n")
		req.InvocContext().GetNode()
		iam := &IAMOutput{}
		if len(req.Arguments()) > 0 {
			key := req.Arguments()[0]
			iam.Key = key
			iam.KeyHash = b58.Encode(util.Hash([]byte(key)))
		}

		iam.Ih, _, _ = req.Option("ih").String()
		iam.Iv, _, _ = req.Option("iv").String()

		res.SetOutput(iam)
	},
	Type: IAMOutput{},
}

var IamInfoNewCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM CMD",
		ShortDescription: "iam use new command",
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("key", true, true, "The path to the IPFS object(s) to be outputted."),
	},
	Options: []cmdkit.Option{
		cmdkit.StringOption("ih", "Byte offset to begin reading from."),
		cmdkit.StringOption("iv", "Maximum number of bytes to read."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {
		fmt.Printf("load info ...\n")

		ctx, _ := env.(*oldcmds.Context)
		ctx.GetNode()

		iam := &IAMOutput{}
		if len(req.Arguments) > 0 {
			key := req.Arguments[0]
			iam.Key = key
			iam.KeyHash = b58.Encode(util.Hash([]byte(key)))
		}
		iam.Ih, _ = req.Options["ih"].(string)
		iam.Iv, _ = req.Options["iv"].(string)
		cmds.EmitOnce(res, &iam)
	},
	Type: IAMOutput{},
}
