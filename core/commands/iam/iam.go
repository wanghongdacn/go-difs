package iam

import (
	"fmt"
	oldcmds "github.com/Harold-the-Axeman/dacc-iam-filesystem/commands"
	lgc "github.com/Harold-the-Axeman/dacc-iam-filesystem/commands/legacy"
	b58 "github.com/mr-tron/base58/base58"
	cmds "gx/ipfs/QmPTfgFTo9PFr1PvPKyKoeMgBvYPh6cX3aDP7DHKVbnCbi/go-ipfs-cmds"
	"gx/ipfs/QmPdKqUcHGFdeSpvjVoaTRPPstGif9GBZb5Q56RVw9o69A/go-ipfs-util"
	"gx/ipfs/QmSP88ryZkHSRn1fnngAaV2Vcn63WUJzAavnRM9CVdU1Ky/go-ipfs-cmdkit"
)

type IAMOutput struct {
	ID        string
	PublicKey string
}

var IamCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM CMD",
		ShortDescription: "",
	},
	Subcommands: map[string]*cmds.Command{
		"info": lgc.NewCommand(IamInfoCmd),
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
		cmdkit.StringArg("key", true, true, "The path to the IPFS object(s) to be outputted.").EnableStdin(),
	},
	Options: []cmdkit.Option{
		cmdkit.IntOption("ih", "Byte offset to begin reading from."),
		cmdkit.IntOption("iv", "Maximum number of bytes to read."),
	},
	Run: func(req oldcmds.Request, res oldcmds.Response) {
		fmt.Printf("load info ...\n")

		var id string
		var pk string
		if len(req.Arguments()) > 0 {
			var err error
			id = req.Arguments()[0]
			pk = b58.Encode(util.Hash([]byte(id)))
			if err != nil {
				res.SetError(cmds.ClientError("Invalid id"), cmdkit.ErrClient)
				return
			}
		} else {
			id = "none"
			pk = "none"
		}

		iam := &IAMOutput{
			ID:        id,
			PublicKey: pk,
		}
		res.SetOutput(iam)
	},
	Type: IAMOutput{},
}
