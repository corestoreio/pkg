// +build ignore

package catalogurlrewrite

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCatalogSeoCategoryUrlSuffix => Category URL Suffix.
// You need to refresh the cache.
// BackendModel: Otnegam\Catalog\Model\System\Config\Backend\Catalog\Url\Rewrite\Suffix
var PathCatalogSeoCategoryUrlSuffix = model.NewStr(`catalog/seo/category_url_suffix`)

// PathCatalogSeoProductUrlSuffix => Product URL Suffix.
// You need to refresh the cache.
// BackendModel: Otnegam\Catalog\Model\System\Config\Backend\Catalog\Url\Rewrite\Suffix
var PathCatalogSeoProductUrlSuffix = model.NewStr(`catalog/seo/product_url_suffix`)

// PathCatalogSeoProductUseCategories => Use Categories Path for Product URLs.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogSeoProductUseCategories = model.NewBool(`catalog/seo/product_use_categories`)

// PathCatalogSeoSaveRewritesHistory => Create Permanent Redirect for URLs if URL Key Changed.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogSeoSaveRewritesHistory = model.NewBool(`catalog/seo/save_rewrites_history`)
