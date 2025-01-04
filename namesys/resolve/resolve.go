package resolve

import (
	"context"
	"errors"
	"fmt"
	"strings"

	logging "github.com/ipfs/go-log"
	"github.com/ipfs/go-path"

	"github.com/ipfs/go-ipfs/namesys"
)

var log = logging.Logger("nsresolv")

// ErrNoNamesys is an explicit error for when an IPFS node doesn't
// (yet) have a name system
var ErrNoNamesys = errors.New(
	"core/resolve: no Namesys on IpfsNode - can't resolve ipns entry")

// ResolveIPNS resolves /ipns paths
func ResolveIPNS(ctx context.Context, nsys namesys.NameSystem, p path.Path) (path.Path, error) {
	if strings.HasPrefix(p.String(), "/ipns/") {
		evt := log.EventBegin(ctx, "resolveIpnsPath")
		defer evt.Done()
		// resolve ipns paths

		// TODO(cryptix): we should be able to query the local cache for the path
		if nsys == nil {
			evt.Append(logging.LoggableMap{"error": ErrNoNamesys.Error()})
			return "", ErrNoNamesys
		}

		seg := p.Segments()

		if len(seg) < 2 || seg[1] == "" { // just "/<protocol/>" without further segments
			err := fmt.Errorf("invalid path %q: ipns path missing IPNS ID", p)
			evt.Append(logging.LoggableMap{"error": err})
			return "", err
		}

		extensions := seg[2:]
		resolvable, err := path.FromSegments("/", seg[0], seg[1])
		if err != nil {
			evt.Append(logging.LoggableMap{"error": err.Error()})
			return "", err
		}

		respath, err := nsys.Resolve(ctx, resolvable.String())
		if err != nil {
			evt.Append(logging.LoggableMap{"error": err.Error()})
			return "", err
		}

		segments := append(respath.Segments(), extensions...)
		p, err = path.FromSegments("/", segments...)
		if err != nil {
			evt.Append(logging.LoggableMap{"error": err.Error()})
			return "", err
		}
	}
	return p, nil
}
