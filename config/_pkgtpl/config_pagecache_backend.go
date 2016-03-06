// +build ignore

package pagecache

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
	// SystemFullPageCacheCachingApplication => Caching Application.
	// Path: system/full_page_cache/caching_application
	// SourceModel: Magento\PageCache\Model\System\Config\Source\Application
	SystemFullPageCacheCachingApplication model.Str

	// SystemFullPageCacheTtl => TTL for public content.
	// Public content cache lifetime in seconds. If field is empty default value
	// 86400 will be saved.
	// Path: system/full_page_cache/ttl
	// BackendModel: Magento\PageCache\Model\System\Config\Backend\Ttl
	SystemFullPageCacheTtl model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.SystemFullPageCacheCachingApplication = model.NewStr(`system/full_page_cache/caching_application`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SystemFullPageCacheTtl = model.NewStr(`system/full_page_cache/ttl`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
