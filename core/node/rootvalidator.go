package node

import (
	"context"
	"fmt"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/libp2p/go-libp2p-kad-dht/analysis"
	boom "github.com/tylertreat/BoomFilters"
	"time"
)

type rootValidator struct {
	rootValidator analysis.RootValidator
	blockStore    blockstore.Blockstore
	// This is a bloom-like structure that contains all the non-root block ids
	rootStore        *boom.CuckooFilter
	rootSearchTicker *time.Ticker
	dag              ipld.DAGService
}

const (
	DefaultExpectedCIDs = 100000
)

func (s *rootValidator) FindRoots(bs blockstore.Blockstore, dag ipld.DAGService) error {

	context, _ := context.WithTimeout(context.Background(), 60*time.Second)

	blks, _ := bs.AllKeysChan(context)

	s.rootStore.Reset()

	for blk := range blks {
		n, err := dag.Get(context, blk)

		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, lnk := range n.Links() {
			s.rootStore.Add(lnk.Cid.Bytes())
		}
	}

	return nil
}

func (s *rootValidator) CheckRoot(cid cid.Cid) (bool, error) {

	return !s.rootStore.Test(cid.Bytes()), nil
}

func NewRootValidator(bs blockstore.Blockstore, dag ipld.DAGService) analysis.RootValidator {

	ticker := time.NewTicker(5 * time.Minute)

	v := &rootValidator{
		blockStore:       bs,
		dag:              dag,
		rootStore:        boom.NewCuckooFilter(DefaultExpectedCIDs, 0.05),
		rootSearchTicker: ticker,
	}

	//	quit := make(chan struct{})
	go func() {
		// Do a first search immediately
		v.FindRoots(v.blockStore, v.dag)

		for {
			select {
			case <-ticker.C:
				v.FindRoots(v.blockStore, v.dag)

			}
			//case <- quit:
			//	ticker.Stop()
			//	return
			//}
		}
	}()

	return v
}
