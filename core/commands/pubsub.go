package commands

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"

	core "github.com/ipfs/go-ipfs/core"
	cmdenv "github.com/ipfs/go-ipfs/core/commands/cmdenv"
	e "github.com/ipfs/go-ipfs/core/commands/e"

	cid "gx/ipfs/QmPSQnBKM9g7BaUcZCvswUJVscQ1ipjmwxN5PXCjkp9EQ7/go-cid"
	cmds "gx/ipfs/QmPXR4tNdLbp8HsZiPMjpsgqphX9Vhw2J6Jh5MKH2ovW3D/go-ipfs-cmds"
	blocks "gx/ipfs/QmRcHuYzAyswytBuMF78rj3LTChYszomRFXNg4685ZN1WM/go-block-format"
	cmdkit "gx/ipfs/QmSP88ryZkHSRn1fnngAaV2Vcn63WUJzAavnRM9CVdU1Ky/go-ipfs-cmdkit"
	floodsub "gx/ipfs/QmY1L5krVk8dv8d74uESmJTXGpoigVYqBVxXXz1aS8aFSb/go-libp2p-floodsub"
	pstore "gx/ipfs/Qmda4cPRvSRyox3SqgJN6DfSZGU5TtHufPTp9uXjFj71X6/go-libp2p-peerstore"
)

var PubsubCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "An experimental publish-subscribe system on ipfs.",
		ShortDescription: `
ipfs pubsub allows you to publish messages to a given topic, and also to
subscribe to new messages on a given topic.

This is an experimental feature. It is not intended in its current state
to be used in a production environment.

To use, the daemon must be run with '--enable-pubsub-experiment'.
`,
	},
	Subcommands: map[string]*cmds.Command{
		"pub":   PubsubPubCmd,
		"sub":   PubsubSubCmd,
		"ls":    PubsubLsCmd,
		"peers": PubsubPeersCmd,
	},
}

var PubsubSubCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Subscribe to messages on a given topic.",
		ShortDescription: `
ipfs pubsub sub subscribes to messages on a given topic.

This is an experimental feature. It is not intended in its current state
to be used in a production environment.

To use, the daemon must be run with '--enable-pubsub-experiment'.
`,
		LongDescription: `
ipfs pubsub sub subscribes to messages on a given topic.

This is an experimental feature. It is not intended in its current state
to be used in a production environment.

To use, the daemon must be run with '--enable-pubsub-experiment'.

This command outputs data in the following encodings:
  * "json"
(Specified by the "--encoding" or "--enc" flag)
`,
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("topic", true, false, "String name of topic to subscribe to."),
	},
	Options: []cmdkit.Option{
		cmdkit.BoolOption("discover", "try to discover other peers subscribed to the same topic"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		// Must be online!
		if !n.OnlineMode() {
			return cmdkit.Errorf(cmdkit.ErrClient, ErrNotOnline.Error())
		}

		if n.Floodsub == nil {
			return fmt.Errorf("experimental pubsub feature not enabled. Run daemon with --enable-pubsub-experiment to use")
		}

		topic := req.Arguments[0]
		sub, err := n.Floodsub.Subscribe(topic)
		if err != nil {
			return err
		}
		defer sub.Cancel()

		discover, _ := req.Options["discover"].(bool)
		if discover {
			go func() {
				blk := blocks.NewBlock([]byte("floodsub:" + topic))
				err := n.Blocks.AddBlock(blk)
				if err != nil {
					log.Error("pubsub discovery: ", err)
					return
				}

				connectToPubSubPeers(req.Context, n, blk.Cid())
			}()
		}

		if f, ok := res.(http.Flusher); ok {
			f.Flush()
		}

		for {
			msg, err := sub.Next(req.Context)
			if err == io.EOF || err == context.Canceled {
				return nil
			} else if err != nil {
				return err
			}

			err = res.Emit(msg)
			if err != nil {
				return err
			}
		}
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeEncoder(func(req *cmds.Request, w io.Writer, v interface{}) error {
			m, ok := v.(*floodsub.Message)
			if !ok {
				return fmt.Errorf("unexpected type: %T", v)
			}

			_, err := w.Write(m.Data)
			return err
		}),
		"ndpayload": cmds.MakeEncoder(func(req *cmds.Request, w io.Writer, v interface{}) error {
			m, ok := v.(*floodsub.Message)
			if !ok {
				return fmt.Errorf("unexpected type: %T", v)
			}

			m.Data = append(m.Data, '\n')
			_, err := w.Write(m.Data)
			return err
		}),
		"lenpayload": cmds.MakeEncoder(func(req *cmds.Request, w io.Writer, v interface{}) error {
			m, ok := v.(*floodsub.Message)
			if !ok {
				return fmt.Errorf("unexpected type: %T", v)
			}

			buf := make([]byte, 8, len(m.Data)+8)

			n := binary.PutUvarint(buf, uint64(len(m.Data)))
			buf = append(buf[:n], m.Data...)
			_, err := w.Write(buf)
			return err
		}),
	},
	Type: floodsub.Message{},
}

func connectToPubSubPeers(ctx context.Context, n *core.IpfsNode, cid cid.Cid) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	provs := n.Routing.FindProvidersAsync(ctx, cid.Hash(), 10)
	wg := &sync.WaitGroup{}
	for p := range provs {
		wg.Add(1)
		go func(pi pstore.PeerInfo) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(ctx, time.Second*10)
			defer cancel()
			err := n.PeerHost.Connect(ctx, pi)
			if err != nil {
				log.Info("pubsub discover: ", err)
				return
			}
			log.Info("connected to pubsub peer:", pi.ID)
		}(p)
	}

	wg.Wait()
}

var PubsubPubCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Publish a message to a given pubsub topic.",
		ShortDescription: `
ipfs pubsub pub publishes a message to a specified topic.

This is an experimental feature. It is not intended in its current state
to be used in a production environment.

To use, the daemon must be run with '--enable-pubsub-experiment'.
`,
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("topic", true, false, "Topic to publish to."),
		cmdkit.StringArg("data", true, true, "Payload of message to publish.").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		// Must be online!
		if !n.OnlineMode() {
			return cmdkit.Errorf(cmdkit.ErrClient, ErrNotOnline.Error())
		}

		if n.Floodsub == nil {
			return errors.New("experimental pubsub feature not enabled. Run daemon with --enable-pubsub-experiment to use.")
		}

		topic := req.Arguments[0]

		err = req.ParseBodyArgs()
		if err != nil {
			return err
		}

		for _, data := range req.Arguments[1:] {
			if err := n.Floodsub.Publish(topic, []byte(data)); err != nil {
				return err
			}
		}

		return nil
	},
}

var PubsubLsCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "List subscribed topics by name.",
		ShortDescription: `
ipfs pubsub ls lists out the names of topics you are currently subscribed to.

This is an experimental feature. It is not intended in its current state
to be used in a production environment.

To use, the daemon must be run with '--enable-pubsub-experiment'.
`,
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		// Must be online!
		if !n.OnlineMode() {
			return cmdkit.Errorf(cmdkit.ErrClient, ErrNotOnline.Error())
		}

		if n.Floodsub == nil {
			return errors.New("experimental pubsub feature not enabled. Run daemon with --enable-pubsub-experiment to use.")
		}

		return cmds.EmitOnce(res, stringList{n.Floodsub.GetTopics()})
	},
	Type: stringList{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeEncoder(stringListEncoder),
	},
}

func stringListEncoder(req *cmds.Request, w io.Writer, v interface{}) error {
	list, ok := v.(*stringList)
	if !ok {
		return e.TypeErr(list, v)
	}
	for _, str := range list.Strings {
		_, err := fmt.Fprintf(w, "%s\n", str)
		if err != nil {
			return err
		}
	}
	return nil
}

var PubsubPeersCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "List peers we are currently pubsubbing with.",
		ShortDescription: `
ipfs pubsub peers with no arguments lists out the pubsub peers you are
currently connected to. If given a topic, it will list connected
peers who are subscribed to the named topic.

This is an experimental feature. It is not intended in its current state
to be used in a production environment.

To use, the daemon must be run with '--enable-pubsub-experiment'.
`,
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("topic", false, false, "topic to list connected peers of"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		n, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		// Must be online!
		if !n.OnlineMode() {
			return cmdkit.Errorf(cmdkit.ErrClient, ErrNotOnline.Error())
		}

		if n.Floodsub == nil {
			return errors.New("experimental pubsub feature not enabled. Run daemon with --enable-pubsub-experiment to use")
		}

		var topic string
		if len(req.Arguments) == 1 {
			topic = req.Arguments[0]
		}

		peers := n.Floodsub.ListPeers(topic)
		list := &stringList{make([]string, 0, len(peers))}

		for _, peer := range peers {
			list.Strings = append(list.Strings, peer.Pretty())
		}
		sort.Strings(list.Strings)
		return cmds.EmitOnce(res, list)
	},
	Type: stringList{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeEncoder(stringListEncoder),
	},
}
