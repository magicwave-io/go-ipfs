package commands

import (
	"io"
	"strings"

	cmds "github.com/ipfs/go-ipfs/commands"
	core "github.com/ipfs/go-ipfs/core"
	cmdenv "github.com/ipfs/go-ipfs/core/commands/cmdenv"
	e "github.com/ipfs/go-ipfs/core/commands/e"
	"github.com/ipfs/go-ipfs/core/coreunix"
	tar "github.com/ipfs/go-ipfs/tar"
	dag "gx/ipfs/QmcBoNcAP6qDjgRBew7yjvCqHq7p5jMstE44jPUBWBxzsV/go-merkledag"
	path "gx/ipfs/QmcjwUb36Z16NJkvDX6ccXPqsFswo6AsRXynyXcLLCphV2/go-path"

	"gx/ipfs/QmSP88ryZkHSRn1fnngAaV2Vcn63WUJzAavnRM9CVdU1Ky/go-ipfs-cmdkit"
	apicid "gx/ipfs/QmdPF1WZQHFNfLdwhaShiR3e4KvFviAM58TrxVxPMhukic/go-cidutil/apicid"
)

var TarCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Utility functions for tar files in ipfs.",
	},

	Subcommands: map[string]*cmds.Command{
		"add": tarAddCmd,
		"cat": tarCatCmd,
	},
}

var tarAddCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Import a tar file into ipfs.",
		ShortDescription: `
'ipfs tar add' will parse a tar file and create a merkledag structure to
represent it.
`,
	},

	Arguments: []cmdkit.Argument{
		cmdkit.FileArg("file", true, false, "Tar file to add.").EnableStdin(),
	},
	Run: func(req cmds.Request, res cmds.Response) {
		nd, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		fi, err := req.Files().NextFile()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		node, err := tar.ImportTar(req.Context(), fi, nd.DAG)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		c := node.Cid()

		fi.FileName()
		res.SetOutput(&coreunix.AddedObject{
			Name: fi.FileName(),
			Hash: apicid.FromCid(c),
		})
	},
	Type: coreunix.AddedObject{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			_, err := cmdenv.NewCidBaseHandlerLegacy(res.Request()).UseGlobal().Proc()
			if err != nil {
				return nil, err
			}

			v, err := unwrapOutput(res.Output())
			if err != nil {
				return nil, err
			}

			o, ok := v.(*coreunix.AddedObject)
			if !ok {
				return nil, e.TypeErr(o, v)
			}
			return strings.NewReader(o.Hash.String() + "\n"), nil
		},
	},
}

var tarCatCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Export a tar file from IPFS.",
		ShortDescription: `
'ipfs tar cat' will export a tar file from a previously imported one in IPFS.
`,
	},

	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("path", true, false, "ipfs path of archive to export.").EnableStdin(),
	},
	Run: func(req cmds.Request, res cmds.Response) {
		nd, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		p, err := path.ParsePath(req.Arguments()[0])
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		root, err := core.Resolve(req.Context(), nd.Namesys, nd.Resolver, p)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		rootpb, ok := root.(*dag.ProtoNode)
		if !ok {
			res.SetError(dag.ErrNotProtobuf, cmdkit.ErrNormal)
			return
		}

		r, err := tar.ExportTar(req.Context(), rootpb, nd.DAG)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		res.SetOutput(r)
	},
}
