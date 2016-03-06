// +build ignore

package integration

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	model.PkgBackend
	// OauthCleanupCleanupProbability => Cleanup Probability.
	// Integer. Launch cleanup in X OAuth requests. 0 (not recommended) - to
	// disable cleanup
	// Path: oauth/cleanup/cleanup_probability
	OauthCleanupCleanupProbability model.Str

	// OauthCleanupExpirationPeriod => Expiration Period.
	// Cleanup entries older than X minutes.
	// Path: oauth/cleanup/expiration_period
	OauthCleanupExpirationPeriod model.Str

	// OauthConsumerExpirationPeriod => Expiration Period.
	// Consumer key/secret will expire if not used within X seconds after Oauth
	// token exchange starts.
	// Path: oauth/consumer/expiration_period
	OauthConsumerExpirationPeriod model.Str

	// OauthConsumerPostMaxredirects => OAuth consumer credentials HTTP Post maxredirects.
	// Number of maximum redirects for OAuth consumer credentials Post request.
	// Path: oauth/consumer/post_maxredirects
	OauthConsumerPostMaxredirects model.Str

	// OauthConsumerPostTimeout => OAuth consumer credentials HTTP Post timeout.
	// Timeout for OAuth consumer credentials Post request within X seconds.
	// Path: oauth/consumer/post_timeout
	OauthConsumerPostTimeout model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.OauthCleanupCleanupProbability = model.NewStr(`oauth/cleanup/cleanup_probability`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.OauthCleanupExpirationPeriod = model.NewStr(`oauth/cleanup/expiration_period`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.OauthConsumerExpirationPeriod = model.NewStr(`oauth/consumer/expiration_period`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.OauthConsumerPostMaxredirects = model.NewStr(`oauth/consumer/post_maxredirects`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.OauthConsumerPostTimeout = model.NewStr(`oauth/consumer/post_timeout`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
