package coreapi

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/ipfs/go-cid"
	"paidpiper.com/payment-gateway/boom/data"
	boomModule "paidpiper.com/payment-gateway/boom/module"
	boomServer "paidpiper.com/payment-gateway/boom/server"
)

type PlusAPI interface {
	Fill(context.Context, string, chan string) error
}

func (api *CoreAPI) Plus() PlusAPI {
	return api
}

func (api *CoreAPI) Fill(ctx context.Context, s string, ch chan string) error {
	p := &ApiProxy{
		api,
		ctx,
		ch,
	}
	return boomModule.Fill(p, p, ch)
}

type ApiProxy struct {
	api *CoreAPI
	ctx context.Context
	ch  chan string
}

func (prox *ApiProxy) Get(b [][]byte) error {
	cids := []cid.Cid{}
	for _, bs := range b {
		cid, err := cid.Parse(bs)
		if err != nil {
			return err
		}
		cids = append(cids, cid)
	}
	blocksChain := prox.api.blocks.GetBlocks(prox.ctx, cids)
	for b := range blocksChain {
		prox.ch <- fmt.Sprintf("Get Block %v %v bytes", b.Cid().String(), len(b.RawData()))

		prox.api.blocks.AddBlock(b)
	}
	return nil
}
func isChanClosed(ch interface{}) bool {
	if reflect.TypeOf(ch).Kind() != reflect.Chan {
		panic("only channels!")
	}

	// get interface value pointer, from cgo_export
	// typedef struct { void *t; void *v; } GoInterface;
	// then get channel real pointer
	cptr := *(*uintptr)(unsafe.Pointer(
		unsafe.Pointer(uintptr(unsafe.Pointer(&ch)) + unsafe.Sizeof(uint(0))),
	))

	// this function will return true if chan.closed > 0
	// see hchan on https://github.com/golang/go/blob/master/src/runtime/chan.go
	// type hchan struct {
	// qcount   uint           // total data in the queue
	// dataqsiz uint           // size of the circular queue
	// buf      unsafe.Pointer // points to an array of dataqsiz elements
	// elemsize uint16
	// closed   uint32
	// **

	cptr += unsafe.Sizeof(uint(0)) * 2
	cptr += unsafe.Sizeof(unsafe.Pointer(uintptr(0)))
	cptr += unsafe.Sizeof(uint16(0))
	return *(*uint32)(unsafe.Pointer(cptr)) > 0
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
			addr := address.String()
			if strings.HasPrefix(addr, "/onion") {
				connections.Hosts = append(connections.Hosts, addr)
			}
		}

	}
	return connections, nil
}
