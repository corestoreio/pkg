// +build ignore

package catalogurlrewrite

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
	// CatalogSeoCategoryUrlSuffix => Category URL Suffix.
	// You need to refresh the cache.
	// Path: catalog/seo/category_url_suffix
	// BackendModel: Otnegam\Catalog\Model\System\Config\Backend\Catalog\Url\Rewrite\Suffix
	CatalogSeoCategoryUrlSuffix model.Str

	// CatalogSeoProductUrlSuffix => Product URL Suffix.
	// You need to refresh the cache.
	// Path: catalog/seo/product_url_suffix
	// BackendModel: Otnegam\Catalog\Model\System\Config\Backend\Catalog\Url\Rewrite\Suffix
	CatalogSeoProductUrlSuffix model.Str

	// CatalogSeoProductUseCategories => Use Categories Path for Product URLs.
	// Path: catalog/seo/product_use_categories
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogSeoProductUseCategories model.Bool

	// CatalogSeoSaveRewritesHistory => Create Permanent Redirect for URLs if URL Key Changed.
	// Path: catalog/seo/save_rewrites_history
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogSeoSaveRewritesHistory model.Bool
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogSeoCategoryUrlSuffix = model.NewStr(`catalog/seo/category_url_suffix`, model.WithPkgCfg(pkgCfg))
	pp.CatalogSeoProductUrlSuffix = model.NewStr(`catalog/seo/product_url_suffix`, model.WithPkgCfg(pkgCfg))
	pp.CatalogSeoProductUseCategories = model.NewBool(`catalog/seo/product_use_categories`, model.WithPkgCfg(pkgCfg))
	pp.CatalogSeoSaveRewritesHistory = model.NewBool(`catalog/seo/save_rewrites_history`, model.WithPkgCfg(pkgCfg))

	return pp
}
