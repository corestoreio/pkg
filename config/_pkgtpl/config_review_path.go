// +build ignore

package review

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
	// CatalogReviewAllowGuest => Allow Guests to Write Reviews.
	// Path: catalog/review/allow_guest
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogReviewAllowGuest model.Bool
}

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogReviewAllowGuest = model.NewBool(`catalog/review/allow_guest`, model.WithConfigStructure(cfgStruct))

	return pp
}
