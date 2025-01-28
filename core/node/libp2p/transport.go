package libp2p

import (
	"fmt"
	oniontransport "paidpiper.com/go-libp2p-onion-transport"

	config "github.com/ipfs/go-ipfs-config"
	libp2p "github.com/libp2p/go-libp2p"
	metrics "github.com/libp2p/go-libp2p-core/metrics"
	libp2pquic "github.com/libp2p/go-libp2p-quic-transport"
	tcp "github.com/libp2p/go-tcp-transport"
	websocket "github.com/libp2p/go-ws-transport"

	"go.uber.org/fx"
)


func TorPath(path string) func() (opts Libp2pOpts, err error) {
	return func() (opts Libp2pOpts, err error) {
		opts.Opts = append(opts.Opts, libp2p.TorPath(path))

		return opts, nil
	}
}

func TorConfigPath(path string) func() (opts Libp2pOpts, err error) {
	return func() (opts Libp2pOpts, err error) {
		opts.Opts = append(opts.Opts, libp2p.TorConfigPath(path))

		return opts, nil
	}
}

func Transports(tptConfig config.Transports, config *config.Config) interface{} {
	return func(pnet struct {
		fx.In
		Fprint PNetFingerprint `optional:"true"`
	}) (opts Libp2pOpts,err error) {
		privateNetworkEnabled := pnet.Fprint != nil



		if tptConfig.Network.TCP.WithDefault(true) {
			opts.Opts = append(opts.Opts, libp2p.Transport(tcp.NewTCPTransport))
		}

		if tptConfig.Network.Websocket.WithDefault(true) {
			opts.Opts = append(opts.Opts, libp2p.Transport(websocket.New))
		}

		if tptConfig.Network.Tor.WithDefault(true) {
			opts.Opts = append(opts.Opts, libp2p.Transport(
				oniontransport.NewOnionTransportC(config.TorPath,
					config.TorConfigPath,"", nil, "", true)))
		}

		if tptConfig.Network.QUIC.WithDefault(!privateNetworkEnabled) {
			if privateNetworkEnabled {
				// QUIC was force enabled while the private network was turned on.
				// Fail and tell the user.
				return opts, fmt.Errorf(
					"The QUIC transport does not support private networks. " +
						"Please disable Swarm.Transports.Network.QUIC.",
				)
			}
			opts.Opts = append(opts.Opts, libp2p.Transport(libp2pquic.NewTransport))
		}

		return opts, nil
	}
}

func BandwidthCounter() (opts Libp2pOpts, reporter *metrics.BandwidthCounter) {
	reporter = metrics.NewBandwidthCounter()
	opts.Opts = append(opts.Opts, libp2p.BandwidthReporter(reporter))
	return opts, reporter
}
