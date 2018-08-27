package commands

import (
	"fmt"
	"gx/ipfs/QmPTfgFTo9PFr1PvPKyKoeMgBvYPh6cX3aDP7DHKVbnCbi/go-ipfs-cmds"
	ft "gx/ipfs/QmQjEpRiwVvtowhq69dAtB4jhioPVFXiCcWZm9Sfgn7eqc/go-unixfs"
	dag "gx/ipfs/QmRiQCJZ91B7VNmLvA6sxzDuBJGSojS3uXHHVuNr3iueNZ/go-merkledag"
	"gx/ipfs/QmSP88ryZkHSRn1fnngAaV2Vcn63WUJzAavnRM9CVdU1Ky/go-ipfs-cmdkit"
	ipld "gx/ipfs/QmX5CsuHyVZeTLxgRSYkgLSDQKb9UjE8xnhQzCEJWWWFsC/go-ipld-format"
)

var IamTestCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM test cmd",
		ShortDescription: "",
	},
	Subcommands: map[string]*cmds.Command{
		"cid": IamTestCidCmd,
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
