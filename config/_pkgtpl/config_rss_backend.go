// +build ignore

package rss

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
	// RssConfigActive => Enable RSS.
	// Path: rss/config/active
	// BackendModel: Magento\Rss\Model\System\Config\Backend\Links
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	RssConfigActive model.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.RssConfigActive = model.NewBool(`rss/config/active`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
