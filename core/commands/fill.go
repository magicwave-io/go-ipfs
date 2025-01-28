package commands

import (
	"fmt"
	"io"

	cmds "github.com/ipfs/go-ipfs-cmds"
	"github.com/ipfs/go-ipfs/core/commands/cmdenv"
	"github.com/ipfs/go-ipfs/core/coreapi"
)

var FillCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List the logging subsystems.",
		ShortDescription: `
		`,
	},

	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		api, err := cmdenv.GetApi(env, req)
		if err != nil {
			return err
		}
		apiPlus := api.(*coreapi.CoreAPI)
		err = apiPlus.Plus().Fill(req.Context, "")
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &fullOptions{})
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, list *fullOptions) error {
			fmt.Fprintln(w, list.Size)
			return nil
		}),
	},
	Type: fullOptions{},
}

type fullOptions struct {
	Size string
}
