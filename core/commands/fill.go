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
		ch := make(chan interface{})
		outData := make(chan string)
		chanClosed := false
		go func() {

			go func() {
				for item := range outData {
					if !chanClosed {
						ch <- &fullOptions{
							Description: item,
						}
					} else {
						fmt.Println(item)
					}

				}
			}()
			err := apiPlus.Plus().Fill(req.Context, "", outData)
			if err != nil {
				ch <- err
				return
			}

			chanClosed = true
			close(outData)
			close(ch)
		}()
		return cmds.EmitChan(res, ch)

	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, list *fullOptions) error {
			fmt.Fprintln(w, fmt.Sprintf("%v %v", list.Step, list.Description))
			return nil
		}),
	},
	Type: fullOptions{},
}

type fullOptions struct {
	Step        string
	Description string
}
