package memstore

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/net/ratelimit"
	"github.com/corestoreio/csfw/net/ratelimit/backendratelimit"
	"github.com/corestoreio/csfw/util/errors"
)

func RegisterOptionFacory(be *backendratelimit.Backend) (string, ratelimit.OptionFactoryFunc) {
	return "memstore", func(sg config.ScopedGetter) []ratelimit.Option {

		burst, _, err := be.RateLimitBurst.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[memstore] RateLimitBurst.Get"))
		}
		req, _, err := be.RateLimitRequests.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[memstore] RateLimitRequests.Get"))
		}
		durRaw, _, err := be.RateLimitDuration.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[memstore] RateLimitDuration.Get"))
		}

		if len(durRaw) != 1 {
			return optError(errors.NewFatalf("[memstore] RateLimitDuration invalid character count: %q. Should be one character long.", durRaw))
		}

		dur := rune(durRaw[0])

		useInMemMaxKeys, scpHash, err := be.RateLimitStorageGcraMaxMemoryKeys.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[memstore] RateLimitStorageGcraMaxMemoryKeys.Get"))
		} else if useInMemMaxKeys > 0 {
			scp, scpID := scpHash.Unpack()
			return WithGCRA(scp, scpID, useInMemMaxKeys, dur, req, burst)
		}
		return optError(errors.NewEmptyf("[memstore] Memstore not active because RateLimitStorageGcraMaxMemoryKeys is %d.", useInMemMaxKeys))
	}
}

func optError(err error) []ratelimit.Option {
	return []ratelimit.Option{func(s *ratelimit.Service) error {
		return err // no need to mask here, not interesting.
	}}
}
