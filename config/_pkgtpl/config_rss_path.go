// +build ignore

package rss

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with ConfigStructure.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// RssConfigActive => Enable RSS.
	// Path: rss/config/active
	// BackendModel: Otnegam\Rss\Model\System\Config\Backend\Links
	// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
	RssConfigActive model.Bool
}

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.RssConfigActive = model.NewBool(`rss/config/active`, model.WithConfigStructure(cfgStruct))

	return pp
}
