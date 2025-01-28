package commands

import (
	"fmt"
	"io"
	"sync"

	cmds "github.com/ipfs/go-ipfs-cmds"
	"github.com/ipfs/go-ipfs/core/commands/cmdenv"
	"github.com/ipfs/go-ipfs/core/coreapi"
)

const (
	PROXY_ADDR = "127.0.0.1:29050"
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

		config, err := cmdenv.GetConfig(env)
		if err != nil {
			return err
		}

		// cfgLocation := ""
		// if cfgLocation != "" {
		// 	if conf, err = cserial.Load(cfgLocation); err != nil {
		// 		return err
		// 	}

		// }

		proxyAdress := PROXY_ADDR
		apiPlus := api.(*coreapi.CoreAPI)

		if config.TorProxyUrl != "" {
			proxyAdress = config.TorProxyUrl
		}

		ch := make(chan interface{})
		outData := make(chan string)
		chanClosed := false
		closeLock := sync.Mutex{}

		go func() {

			go func() {
				for item := range outData {
					func() {
						closeLock.Lock()
						defer closeLock.Unlock()

						if !chanClosed {
							ch <- &fullOptions{
								Description: item,
							}
						} else {
							fmt.Println(item)
						}
					}()

				}
			}()

			err := apiPlus.Plus().Fill(req.Context, "", proxyAdress, outData, config)
			if err != nil {
				ch <- err
				return
			}

			closeLock.Lock()
			defer closeLock.Unlock()
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
