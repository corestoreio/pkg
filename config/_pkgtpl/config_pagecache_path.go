// +build ignore

package pagecache

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathSystemFullPageCacheCachingApplication => Caching Application.
// SourceModel: Otnegam\PageCache\Model\System\Config\Source\Application
var PathSystemFullPageCacheCachingApplication = model.NewStr(`system/full_page_cache/caching_application`)

// PathSystemFullPageCacheTtl => TTL for public content.
// Public content cache lifetime in seconds. If field is empty default value
// 86400 will be saved.
// BackendModel: Otnegam\PageCache\Model\System\Config\Backend\Ttl
var PathSystemFullPageCacheTtl = model.NewStr(`system/full_page_cache/ttl`)
