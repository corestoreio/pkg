// +build ignore

package integration

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with ConfigStructure.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
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

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.OauthCleanupCleanupProbability = model.NewStr(`oauth/cleanup/cleanup_probability`, model.WithConfigStructure(cfgStruct))
	pp.OauthCleanupExpirationPeriod = model.NewStr(`oauth/cleanup/expiration_period`, model.WithConfigStructure(cfgStruct))
	pp.OauthConsumerExpirationPeriod = model.NewStr(`oauth/consumer/expiration_period`, model.WithConfigStructure(cfgStruct))
	pp.OauthConsumerPostMaxredirects = model.NewStr(`oauth/consumer/post_maxredirects`, model.WithConfigStructure(cfgStruct))
	pp.OauthConsumerPostTimeout = model.NewStr(`oauth/consumer/post_timeout`, model.WithConfigStructure(cfgStruct))

	return pp
}
