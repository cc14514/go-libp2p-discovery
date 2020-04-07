package discovery

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"

	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("discovery")

// FindPeers is a utility function that synchronously collects peers from a Discoverer.
func FindPeers(ctx context.Context, d Discoverer, ns string, opts ...Option) ([]peer.AddrInfo, error) {
	var res []peer.AddrInfo

	ch, err := d.FindPeers(ctx, ns, opts...)
	if err != nil {
		return nil, err
	}

	for pi := range ch {
		res = append(res, pi)
	}

	return res, nil
}

// Advertise is a utility function that persistently advertises a service through an Advertiser.
func Advertise(ctx context.Context, a Advertiser, ns string, opts ...Option) {
	go func() {
		for {
			ttl, err := a.Advertise(ctx, ns, opts...)
			log.Infof("discovery::Advertise-1 : ns=%s , ttl=%v , err=%v", ns, ttl, err)
			if err != nil {
				log.Debugf("Error advertising %s: %s", ns, err.Error())
				if ctx.Err() != nil {
					return
				}

				select {
				case <-time.After(2 * time.Minute):
					continue
				case <-ctx.Done():
					return
				}
			}

			wait := 7 * ttl / 8
			log.Infof("discovery::Advertise-2: ns=%s , ttl=%v , wait=%v", ns, ttl, wait)
			select {
			case <-time.After(wait):
				log.Infof("discovery::Advertise-3 loop: ns=%s , ttl=%v , wait=%v", ns, ttl, wait)
			case <-ctx.Done():
				log.Infof("discovery::Advertise-4 done: ns=%s , ttl=%v , wait=%v", ns, ttl, wait)
				return
			}
		}
	}()
}
