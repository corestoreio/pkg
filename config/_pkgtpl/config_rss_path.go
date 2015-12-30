// +build ignore

package rss

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathRssConfigActive => Enable RSS.
// BackendModel: Otnegam\Rss\Model\System\Config\Backend\Links
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathRssConfigActive = model.NewBool(`rss/config/active`)
