package coreapi

import (
	"context"
	"fmt"
	config "github.com/ipfs/go-ipfs-config"
	"reflect"
	"strings"
	"time"
	"unsafe"

	"github.com/ipfs/go-cid"
	"golang.org/x/time/rate"
	"paidpiper.com/payment-gateway/boom/data"
	boomModule "paidpiper.com/payment-gateway/boom/module"
	boomServer "paidpiper.com/payment-gateway/boom/server"
)

type PlusAPI interface {
	Fill(ctx context.Context, s string, proxyPort string, ch chan string, config *config.Config) error
}

func (api *CoreAPI) Plus() PlusAPI {
	return api
}

func (api *CoreAPI) Fill(ctx context.Context, s string, proxyPort string, ch chan string, config *config.Config) error {
	p := &ApiProxy{
		api,
		ctx,
		ch,
		config,
	}
	return boomModule.Fill(p, p, proxyPort, ch)
}

type ApiProxy struct {
	api    *CoreAPI
	ctx    context.Context
	ch     chan string
	config *config.Config
}

func (prox *ApiProxy) Get(b [][]byte) error {

	chunkSize := prox.config.FillChunkSize
	var chunkRetrievalTimeout int32 = 0
	chunkRetrievalTimeout = int32(prox.config.FillChunkRetrievalTimeoutSec)

	// Set to default chunk size if not provided
	if chunkSize == 0 {
		chunkSize = 25
	}

	if chunkRetrievalTimeout == 0 {
		chunkRetrievalTimeout = 30
	}

	// Override 200/60 gives good result
	chunkSize = 10
	chunkRetrievalTimeout = 3
	getTimeout := 3

	//limiterDuration := math.Round(float64(1.1 * chunkRetrievalTimeout))
	t1 := time.Duration(float64(chunkRetrievalTimeout) * float64(time.Second))
	limiter := rate.NewLimiter(rate.Every(t1), 2)
	_ = limiter

	allCids := []cid.Cid{}
	for _, bs := range b {
		cid, err := cid.Parse(bs)
		if err != nil {
			return err
		}
		allCids = append(allCids, cid)
	}
	cidChunk := []cid.Cid{}

	for _, c := range allCids {
		cidChunk = append(cidChunk, c)

		if len(cidChunk) > chunkSize {
			//limiter.Wait(prox.ctx)

			t2 := time.Duration(float64(getTimeout) * float64(time.Second))
			ctx, _ := context.WithTimeout(prox.ctx, t2)

			blocksChain := prox.api.blocks.GetBlocks(ctx, cidChunk)

			for b := range blocksChain {
				prox.ch <- fmt.Sprintf("Get Block %v %v bytes", b.Cid().String(), len(b.RawData()))
				prox.api.blocks.AddBlock(b)
			}

			// Empty the array
			cidChunk = []cid.Cid{}
		}

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
