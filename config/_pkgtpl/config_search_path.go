// +build ignore

package search

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCatalogSearchEngine => Search Engine.
// SourceModel: Otnegam\Search\Model\Adminhtml\System\Config\Source\Engine
var PathCatalogSearchEngine = model.NewStr(`catalog/search/engine`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogSearchSearchType => .
var PathCatalogSearchSearchType = model.NewStr(`catalog/search/search_type`, model.WithPkgCfg(PackageConfiguration))
