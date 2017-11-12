// +build ignore

package catalogurlrewrite

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
	// CatalogSeoCategoryUrlSuffix => Category URL Suffix.
	// You need to refresh the cache.
	// Path: catalog/seo/category_url_suffix
	// BackendModel: Magento\Catalog\Model\System\Config\Backend\Catalog\Url\Rewrite\Suffix
	CatalogSeoCategoryUrlSuffix cfgmodel.Str

	// CatalogSeoProductUrlSuffix => Product URL Suffix.
	// You need to refresh the cache.
	// Path: catalog/seo/product_url_suffix
	// BackendModel: Magento\Catalog\Model\System\Config\Backend\Catalog\Url\Rewrite\Suffix
	CatalogSeoProductUrlSuffix cfgmodel.Str

	// CatalogSeoProductUseCategories => Use Categories Path for Product URLs.
	// Path: catalog/seo/product_use_categories
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogSeoProductUseCategories cfgmodel.Bool

	// CatalogSeoSaveRewritesHistory => Create Permanent Redirect for URLs if URL Key Changed.
	// Path: catalog/seo/save_rewrites_history
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogSeoSaveRewritesHistory cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogSeoCategoryUrlSuffix = cfgmodel.NewStr(`catalog/seo/category_url_suffix`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogSeoProductUrlSuffix = cfgmodel.NewStr(`catalog/seo/product_url_suffix`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogSeoProductUseCategories = cfgmodel.NewBool(`catalog/seo/product_use_categories`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogSeoSaveRewritesHistory = cfgmodel.NewBool(`catalog/seo/save_rewrites_history`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
