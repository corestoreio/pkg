// +build ignore

package integration

import (
	"github.com/corestoreio/pkg/config/cfgmodel"
	"github.com/corestoreio/pkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// OauthCleanupCleanupProbability => Cleanup Probability.
	// Integer. Launch cleanup in X OAuth requests. 0 (not recommended) - to
	// disable cleanup
	// Path: oauth/cleanup/cleanup_probability
	OauthCleanupCleanupProbability cfgmodel.Str

	// OauthCleanupExpirationPeriod => Expiration Period.
	// Cleanup entries older than X minutes.
	// Path: oauth/cleanup/expiration_period
	OauthCleanupExpirationPeriod cfgmodel.Str

	// OauthConsumerExpirationPeriod => Expiration Period.
	// Consumer key/secret will expire if not used within X seconds after Oauth
	// token exchange starts.
	// Path: oauth/consumer/expiration_period
	OauthConsumerExpirationPeriod cfgmodel.Str

	// OauthConsumerPostMaxredirects => OAuth consumer credentials HTTP Post maxredirects.
	// Number of maximum redirects for OAuth consumer credentials Post request.
	// Path: oauth/consumer/post_maxredirects
	OauthConsumerPostMaxredirects cfgmodel.Str

	// OauthConsumerPostTimeout => OAuth consumer credentials HTTP Post timeout.
	// Timeout for OAuth consumer credentials Post request within X seconds.
	// Path: oauth/consumer/post_timeout
	OauthConsumerPostTimeout cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.OauthCleanupCleanupProbability = cfgmodel.NewStr(`oauth/cleanup/cleanup_probability`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.OauthCleanupExpirationPeriod = cfgmodel.NewStr(`oauth/cleanup/expiration_period`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.OauthConsumerExpirationPeriod = cfgmodel.NewStr(`oauth/consumer/expiration_period`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.OauthConsumerPostMaxredirects = cfgmodel.NewStr(`oauth/consumer/post_maxredirects`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.OauthConsumerPostTimeout = cfgmodel.NewStr(`oauth/consumer/post_timeout`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
