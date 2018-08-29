package commands

import (
	"fmt"
	"gx/ipfs/QmPTfgFTo9PFr1PvPKyKoeMgBvYPh6cX3aDP7DHKVbnCbi/go-ipfs-cmds"
	ft "gx/ipfs/QmQjEpRiwVvtowhq69dAtB4jhioPVFXiCcWZm9Sfgn7eqc/go-unixfs"
	dag "gx/ipfs/QmRiQCJZ91B7VNmLvA6sxzDuBJGSojS3uXHHVuNr3iueNZ/go-merkledag"
	"gx/ipfs/QmSP88ryZkHSRn1fnngAaV2Vcn63WUJzAavnRM9CVdU1Ky/go-ipfs-cmdkit"
	ds "gx/ipfs/QmVG5gxteQNEMhrS8prJSmU2C9rebtFuTd3SYZ5kE3YZ5k/go-datastore"
	ipld "gx/ipfs/QmX5CsuHyVZeTLxgRSYkgLSDQKb9UjE8xnhQzCEJWWWFsC/go-ipld-format"
	base32 "gx/ipfs/QmfVj3x4D6Jkq9SEoi5n2NmoUomLwoeiwnYz2KQa15wRw6/base32"
)

var IamTestCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "IAM test cmd",
		ShortDescription: "",
	},
	Subcommands: map[string]*cmds.Command{
		"tt":  IamAutoRunCmd,
		"cid": IamTestCidCmd,
	},
}

var IamAutoRunCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "",
		ShortDescription: "",
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {

		n, _ := GetNode(env)
		fmt.Printf("iam auto run test %s\n", string(n.Identity))

		key := ds.NewKey(base32.RawStdEncoding.EncodeToString([]byte("/pk/QmQ4PSDphDBpFcRRhK9E87Qi2PyJGuXGeuLLnvpBt5Tavg")))
		n.Repo.Datastore().Put(key, []byte("123456"))

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
