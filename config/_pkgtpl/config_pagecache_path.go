// +build ignore

package pagecache

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with PackageConfiguration.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// SystemFullPageCacheCachingApplication => Caching Application.
	// Path: system/full_page_cache/caching_application
	// SourceModel: Otnegam\PageCache\Model\System\Config\Source\Application
	SystemFullPageCacheCachingApplication model.Str

	// SystemFullPageCacheTtl => TTL for public content.
	// Public content cache lifetime in seconds. If field is empty default value
	// 86400 will be saved.
	// Path: system/full_page_cache/ttl
	// BackendModel: Otnegam\PageCache\Model\System\Config\Backend\Ttl
	SystemFullPageCacheTtl model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.SystemFullPageCacheCachingApplication = model.NewStr(`system/full_page_cache/caching_application`, model.WithPkgCfg(pkgCfg))
	pp.SystemFullPageCacheTtl = model.NewStr(`system/full_page_cache/ttl`, model.WithPkgCfg(pkgCfg))

	return pp
}
