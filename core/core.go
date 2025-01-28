/*
Package core implements the IpfsNode object and related methods.

Packages underneath core/ provide a (relatively) stable, low-level API
to carry out most IPFS-related tasks.  For more details on the other
interfaces and how core/... fits into the bigger IPFS picture, see:

  $ godoc github.com/ipfs/go-ipfs
*/
package core

import (
	"context"
	"fmt"
	"github.com/ipfs/go-filestore"
	"github.com/ipfs/go-ipfs-pinner"
	"github.com/libp2p/go-libp2p-kad-dht/analysis"
	"io"

	bserv "github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-graphsync"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
	"github.com/ipfs/go-ipfs-provider"
	ipld "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log"
	mfs "github.com/ipfs/go-mfs"
	resolver "github.com/ipfs/go-path/resolver"
	goprocess "github.com/jbenet/goprocess"
	connmgr "github.com/libp2p/go-libp2p-core/connmgr"
	ic "github.com/libp2p/go-libp2p-core/crypto"
	p2phost "github.com/libp2p/go-libp2p-core/host"
	metrics "github.com/libp2p/go-libp2p-core/metrics"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pstore "github.com/libp2p/go-libp2p-core/peerstore"
	routing "github.com/libp2p/go-libp2p-core/routing"
	ddht "github.com/libp2p/go-libp2p-kad-dht/dual"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	psrouter "github.com/libp2p/go-libp2p-pubsub-router"
	record "github.com/libp2p/go-libp2p-record"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	p2pbhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/ipfs/go-ipfs/core/bootstrap"
	"github.com/ipfs/go-ipfs/core/node"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/fuse/mount"
	"github.com/ipfs/go-ipfs/namesys"
	ipnsrp "github.com/ipfs/go-ipfs/namesys/republisher"
	"github.com/ipfs/go-ipfs/p2p"
	"github.com/ipfs/go-ipfs/peering"
	"github.com/ipfs/go-ipfs/repo"
)

var log = logging.Logger("core")

// IpfsNode is IPFS Core module. It represents an IPFS instance.
type IpfsNode struct {

	// Self
	Identity peer.ID // the local node's identity

	Repo repo.Repo

	// Local node
	Pinning         pin.Pinner             // the pinning manager
	Mounts          Mounts                 `optional:"true"` // current mount state, if any.
	PrivateKey      ic.PrivKey             `optional:"true"` // the local node's private Key
	PNetFingerprint libp2p.PNetFingerprint `optional:"true"` // fingerprint of private network

	// Services
	Peerstore       pstore.Peerstore          `optional:"true"` // storage for other Peer instances
	Blockstore      bstore.GCBlockstore       // the block store (lower level)
	Filestore       *filestore.Filestore      `optional:"true"` // the filestore blockstore
	BaseBlocks      node.BaseBlocks           // the raw blockstore, no filestore wrapping
	GCLocker        bstore.GCLocker           // the locker used to protect the blockstore during gc
	Blocks          bserv.BlockService        // the block service, get/add blocks.
	DAG             ipld.DAGService           // the merkle dag service, get/add objects.
	Resolver        *resolver.Resolver        // the path resolution system
	Reporter        *metrics.BandwidthCounter `optional:"true"`
	Discovery       discovery.Service         `optional:"true"`
	FilesRoot       *mfs.Root
	RecordValidator record.Validator
	RootValidator   analysis.RootValidator

	// Online
	PeerHost      p2phost.Host            `optional:"true"` // the network host (server+client)
	Peering       peering.PeeringService  `optional:"true"`
	Filters       *ma.Filters             `optional:"true"`
	Bootstrapper  io.Closer               `optional:"true"` // the periodic bootstrapper
	Routing       routing.Routing         `optional:"true"` // the routing system. recommend ipfs-dht
	Exchange      exchange.Interface      // the block exchange + strategy (bitswap)
	Namesys       namesys.NameSystem      // the name system, resolves paths to hashes
	Provider      provider.System         // the value provider system
	IpnsRepub     *ipnsrp.Republisher     `optional:"true"`
	GraphExchange graphsync.GraphExchange `optional:"true"`

	PubSub   *pubsub.PubSub             `optional:"true"`
	PSRouter *psrouter.PubsubValueStore `optional:"true"`
	DHT      *ddht.DHT                  `optional:"true"`
	P2P      *p2p.P2P                   `optional:"true"`

	Process goprocess.Process
	ctx     context.Context

	stop func() error

	// Flags
	IsOnline bool `optional:"true"` // Online is set when networking is enabled.
	IsDaemon bool `optional:"true"` // Daemon is set when running on a long-running daemon.
}

// Mounts defines what the node's mount state is. This should
// perhaps be moved to the daemon or mount. It's here because
// it needs to be accessible across daemon requests.
type Mounts struct {
	Ipfs mount.Mount
	Ipns mount.Mount
}

// Close calls Close() on the App object
func (n *IpfsNode) Close() error {
	return n.stop()
}

// Context returns the IpfsNode context
func (n *IpfsNode) Context() context.Context {
	if n.ctx == nil {
		n.ctx = context.TODO()
	}
	return n.ctx
}

type ID string

func splitAddr(m ma.Multiaddr) (transport ma.Multiaddr, id ID) {
	if m == nil {
		return nil, ""
	}

	transport, p2ppart := ma.SplitLast(m)
	if p2ppart == nil || ((p2ppart.Protocol().Code != ma.P_P2P) && (p2ppart.Protocol().Code != ma.P_ONION3)) {
		return m, ""
	}
	id = ID(p2ppart.RawValue()) // already validated by the multiaddr library.
	return transport, id
}

var ErrInvalidAddr = fmt.Errorf("invalid p2p multiaddr")

// PP: Work around to use onion multiaddr
func bootstrapPeers(addrs []string) ([]peer.AddrInfo, error) {

	maddrs := make([]ma.Multiaddr, len(addrs))
	for i, addr := range addrs {
		var err error
		maddrs[i], err = ma.NewMultiaddr(addr)
		if err != nil {
			return nil, err
		}
	}

	//peer.AddrInfosFromP2pAddrs(maddrs...)

	m := make(map[ID][]ma.Multiaddr)
	for _, maddr := range maddrs {
		transport, id := splitAddr(maddr)
		if id == "" {
			return nil, ErrInvalidAddr
		}
		if transport == nil {
			if _, ok := m[id]; !ok {
				m[id] = nil
			}
		} else {
			m[id] = append(m[id], transport)
		}
	}
	ais := make([]peer.AddrInfo, 0, len(m))
	for id, maddrs := range m {
		ais = append(ais, peer.AddrInfo{ID: peer.ID(id), Addrs: maddrs})
	}
	return ais, nil

}

// Bootstrap will set and call the IpfsNodes bootstrap function.
func (n *IpfsNode) Bootstrap(cfg bootstrap.BootstrapConfig) error {
	// TODO what should return value be when in offlineMode?
	if n.Routing == nil {
		return nil
	}

	if n.Bootstrapper != nil {
		n.Bootstrapper.Close() // stop previous bootstrap process.
	}

	// if the caller did not specify a bootstrap peer function, get the
	// freshest bootstrap peers from config. this responds to live changes.
	if cfg.BootstrapPeers == nil {
		cfg.BootstrapPeers = func() []peer.AddrInfo {
			ps, err := n.loadBootstrapPeers()
			if err != nil {
				log.Warn("failed to parse bootstrap peers from config")
				return nil
			}
			return ps
		}
	}

	var err error

	n.Bootstrapper, err = bootstrap.Bootstrap(n.Identity, n.PeerHost, n.Routing, cfg)
	return err
}

func (n *IpfsNode) loadBootstrapPeers() ([]peer.AddrInfo, error) {
	cfg, err := n.Repo.Config()
	if err != nil {
		return nil, err
	}

	return bootstrapPeers(cfg.Bootstrap)
}

type ConstructPeerHostOpts struct {
	AddrsFactory      p2pbhost.AddrsFactory
	DisableNatPortMap bool
	DisableRelay      bool
	EnableRelayHop    bool
	ConnectionManager connmgr.ConnManager
}
