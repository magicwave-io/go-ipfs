module paidpiper/ipfs/go-ipfs

require (
	bazil.org/fuse v0.0.0-20180421153158-65cc252bf669
	github.com/blang/semver v3.5.1+incompatible
	github.com/bren2010/proquint v0.0.0-20160323162903-38337c27106d
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f
	github.com/dustin/go-humanize v1.0.0
	github.com/elgris/jsondiff v0.0.0-20160530203242-765b5c24c302
	github.com/fsnotify/fsnotify v1.4.7
	github.com/go-bindata/go-bindata v3.1.2+incompatible
	github.com/go-delve/delve v1.4.0 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/golang-lru v0.5.3
	github.com/ipfs/go-bitswap v0.1.10
	github.com/ipfs/go-block-format v0.0.2
	github.com/ipfs/go-blockservice v0.1.2
	github.com/ipfs/go-cid v0.0.3
	github.com/ipfs/go-cidutil v0.0.2
	github.com/ipfs/go-datastore v0.3.1
	github.com/ipfs/go-detect-race v0.0.1
	github.com/ipfs/go-ds-badger v0.2.0
	github.com/ipfs/go-ds-flatfs v0.3.0
	github.com/ipfs/go-ds-leveldb v0.4.0
	github.com/ipfs/go-ds-measure v0.1.0
	github.com/ipfs/go-filestore v0.0.3
	github.com/ipfs/go-fs-lock v0.0.3
	github.com/ipfs/go-ipfs v0.4.22
	github.com/ipfs/go-ipfs-blockstore v0.1.0
	github.com/ipfs/go-ipfs-chunker v0.0.3
	github.com/ipfs/go-ipfs-cmds v0.1.1
	github.com/ipfs/go-ipfs-config v0.0.11
	github.com/ipfs/go-ipfs-ds-help v0.0.1
	github.com/ipfs/go-ipfs-exchange-interface v0.0.1
	github.com/ipfs/go-ipfs-exchange-offline v0.0.1
	github.com/ipfs/go-ipfs-files v0.0.4
	github.com/ipfs/go-ipfs-posinfo v0.0.1
	github.com/ipfs/go-ipfs-provider v0.2.2
	github.com/ipfs/go-ipfs-routing v0.1.0
	github.com/ipfs/go-ipfs-util v0.0.1
	github.com/ipfs/go-ipld-cbor v0.0.3
	github.com/ipfs/go-ipld-format v0.0.2
	github.com/ipfs/go-ipld-git v0.0.2
	github.com/ipfs/go-ipns v0.0.1
	github.com/ipfs/go-log v0.0.1
	github.com/ipfs/go-merkledag v0.3.1
	github.com/ipfs/go-metrics-interface v0.0.1
	github.com/ipfs/go-metrics-prometheus v0.0.2
	github.com/ipfs/go-mfs v0.1.1
	github.com/ipfs/go-path v0.0.7
	github.com/ipfs/go-unixfs v0.2.1
	github.com/ipfs/go-verifcid v0.0.1
	github.com/ipfs/interface-go-ipfs-core v0.2.3
	github.com/jbenet/go-is-domain v1.0.3
	github.com/jbenet/go-random v0.0.0-20190219211222-123a90aedc0c
	github.com/jbenet/go-temp-err-catcher v0.0.0-20150120210811-aac704a3f4f2
	github.com/jbenet/goprocess v0.1.3
	github.com/libp2p/go-libp2p v0.4.2
	github.com/libp2p/go-libp2p-autonat-svc v0.1.0
	github.com/libp2p/go-libp2p-circuit v0.1.4
	github.com/libp2p/go-libp2p-connmgr v0.1.1
	github.com/libp2p/go-libp2p-core v0.3.0
	github.com/libp2p/go-libp2p-http v0.1.4
	github.com/libp2p/go-libp2p-kad-dht v0.2.1
	github.com/libp2p/go-libp2p-kbucket v0.2.1
	github.com/libp2p/go-libp2p-loggables v0.1.0
	github.com/libp2p/go-libp2p-mplex v0.2.1
	github.com/libp2p/go-libp2p-peerstore v0.1.4
	github.com/libp2p/go-libp2p-pnet v0.1.0
	github.com/libp2p/go-libp2p-pubsub v0.1.1
	github.com/libp2p/go-libp2p-pubsub-router v0.1.0
	github.com/libp2p/go-libp2p-quic-transport v0.1.1
	github.com/libp2p/go-libp2p-record v0.1.1
	github.com/libp2p/go-libp2p-routing-helpers v0.1.0
	github.com/libp2p/go-libp2p-secio v0.2.1
	github.com/libp2p/go-libp2p-swarm v0.2.2
	github.com/libp2p/go-libp2p-testing v0.1.1
	github.com/libp2p/go-libp2p-tls v0.1.1
	github.com/libp2p/go-libp2p-yamux v0.2.1
	github.com/libp2p/go-maddr-filter v0.0.5
	github.com/libp2p/go-socket-activation v0.0.1
	github.com/libp2p/go-stream-muxer v0.1.0 // indirect
	github.com/libp2p/go-testutil v0.1.0 // indirect
	github.com/miekg/dns v1.1.26 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mr-tron/base58 v1.1.3
	github.com/multiformats/go-multiaddr v0.2.0
	github.com/multiformats/go-multiaddr-dns v0.2.0
	github.com/multiformats/go-multiaddr-net v0.1.1
	github.com/multiformats/go-multibase v0.0.1
	github.com/multiformats/go-multihash v0.0.10
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.1.0
	github.com/stellar/go v0.0.0-20201019170232-24465e0cb5fb
	github.com/syndtr/goleveldb v1.0.0
	github.com/whyrusleeping/base32 v0.0.0-20170828182744-c30ac30633cc
	github.com/whyrusleeping/go-sysinfo v0.0.0-20190219211824-4a357d4b90b1
	github.com/whyrusleeping/multiaddr-filter v0.0.0-20160516205228-e903e4adabd7
	github.com/whyrusleeping/tar-utils v0.0.0-20180509141711-8c6c8ba81d5c
	go.uber.org/fx v1.9.0
	golang.org/x/sys v0.0.0-20191218084908-4a24b4065292
	gopkg.in/cheggaaa/pb.v1 v1.0.28
)

replace github.com/multiformats/go-multiaddr => ../go-multiaddr

replace github.com/multiformats/go-multiaddr-net => ../go-multiaddr-net

replace github.com/libp2p/go-libp2p-core => ../go-libp2p-core

replace paidpiper.com/go-libp2p-onion-transport => ../go-libp2p-onion-transport

replace github.com/libp2p/go-libp2p => ../go-libp2p

replace github.com/libp2p/go-libp2p-swarm => ../go-libp2p-swarm

replace github.com/ipfs/go-bitswap => ../go-bitswap

replace github.com/ipfs/go-ipfs-config => ../go-ipfs-config

replace github.com/xdrpp/stc => github.com/xdrpp/stc v0.0.0-20191113232133-b2821823938c

replace github.com/xdrpp/goxdr => github.com/xdrpp/goxdr v0.0.0-20191113231830-b0b915788678

replace github.com/go-critic/go-critic => github.com/go-critic/go-critic v0.0.0-20190210220443-ee9bf5809ead

replace github.com/golangci/errcheck => github.com/golangci/errcheck v0.0.0-20181223084120-ef45e06d44b6

replace github.com/golangci/go-tools => github.com/golangci/go-tools v0.0.0-20190318060251-af6baa5dc196

replace github.com/golangci/gofmt => github.com/golangci/gofmt v0.0.0-20181222123516-0b8337e80d98

replace github.com/golangci/gosec => github.com/golangci/gosec v0.0.0-20190211064107-66fb7fc33547

replace github.com/golangci/lint-1 => github.com/golangci/lint-1 v0.0.0-20190420132249-ee948d087217

replace mvdan.cc/unparam => mvdan.cc/unparam v0.0.0-20190209190245-fbb59629db34

go 1.14
