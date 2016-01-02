// +build ignore

package catalogsearch

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCatalogSeoSearchTerms => Popular Search Terms.
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathCatalogSeoSearchTerms = model.NewBool(`catalog/seo/search_terms`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogSearchEngine => .
// BackendModel: Otnegam\CatalogSearch\Model\Adminhtml\System\Config\Backend\Engine
var PathCatalogSearchEngine = model.NewStr(`catalog/search/engine`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogSearchMinQueryLength => Minimal Query Length.
var PathCatalogSearchMinQueryLength = model.NewStr(`catalog/search/min_query_length`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogSearchMaxQueryLength => Maximum Query Length.
var PathCatalogSearchMaxQueryLength = model.NewStr(`catalog/search/max_query_length`, model.WithPkgCfg(PackageConfiguration))
