// +build ignore

package integration

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathOauthCleanupCleanupProbability => Cleanup Probability.
// Integer. Launch cleanup in X OAuth requests. 0 (not recommended) - to
// disable cleanup
var PathOauthCleanupCleanupProbability = model.NewStr(`oauth/cleanup/cleanup_probability`, model.WithPkgCfg(PackageConfiguration))

// PathOauthCleanupExpirationPeriod => Expiration Period.
// Cleanup entries older than X minutes.
var PathOauthCleanupExpirationPeriod = model.NewStr(`oauth/cleanup/expiration_period`, model.WithPkgCfg(PackageConfiguration))

// PathOauthConsumerExpirationPeriod => Expiration Period.
// Consumer key/secret will expire if not used within X seconds after Oauth
// token exchange starts.
var PathOauthConsumerExpirationPeriod = model.NewStr(`oauth/consumer/expiration_period`, model.WithPkgCfg(PackageConfiguration))

// PathOauthConsumerPostMaxredirects => OAuth consumer credentials HTTP Post maxredirects.
// Number of maximum redirects for OAuth consumer credentials Post request.
var PathOauthConsumerPostMaxredirects = model.NewStr(`oauth/consumer/post_maxredirects`, model.WithPkgCfg(PackageConfiguration))

// PathOauthConsumerPostTimeout => OAuth consumer credentials HTTP Post timeout.
// Timeout for OAuth consumer credentials Post request within X seconds.
var PathOauthConsumerPostTimeout = model.NewStr(`oauth/consumer/post_timeout`, model.WithPkgCfg(PackageConfiguration))
