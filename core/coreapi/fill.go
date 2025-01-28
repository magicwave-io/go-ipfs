package coreapi

import (
	"context"

	"github.com/ipfs/go-cid"
	"paidpiper.com/payment-gateway/boom/data"
	boomModule "paidpiper.com/payment-gateway/boom/module"
	boomServer "paidpiper.com/payment-gateway/boom/server"
)

type PlusAPI interface {
	Fill(context.Context, string) error
}

func (api *CoreAPI) Plus() PlusAPI {
	return api
}

func (api *CoreAPI) Fill(ctx context.Context, s string) error {
	p := &ApiProxy{
		api,
		ctx,
	}
	return boomModule.Fill(p, p)
}

type ApiProxy struct {
	api *CoreAPI
	ctx context.Context
}

func (prox *ApiProxy) Get(b [][]byte) error {
	cids := []cid.Cid{}
	for _, bs := range b {
		cid, err := cid.Cast(bs)
		if err != nil {
			return err
		}
		cids = append(cids, cid)
	}
	for b := range prox.api.blocks.GetBlocks(prox.ctx, cids) {
		prox.api.blocks.AddBlock(b)
	}
	return nil
}

func (prox *ApiProxy) Elements() ([]*data.FrequencyContentMetadata, error) {
	return boomServer.FrequentElements()
}

func (prox *ApiProxy) Connections() (*data.Connections, error) {
	peerIds := prox.api.peerstore.Peers()
	connections := &data.Connections{}
	for _, peerID := range peerIds {
		peerInfo := prox.api.peerstore.PeerInfo(peerID)
		for _, address := range peerInfo.Addrs {
			connections.Hosts = append(connections.Hosts, address.String())
		}

	}
	return connections, nil
}
