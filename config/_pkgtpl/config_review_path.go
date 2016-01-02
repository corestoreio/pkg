// +build ignore

package review

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
	// CatalogReviewAllowGuest => Allow Guests to Write Reviews.
	// Path: catalog/review/allow_guest
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogReviewAllowGuest model.Bool
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogReviewAllowGuest = model.NewBool(`catalog/review/allow_guest`, model.WithPkgCfg(pkgCfg))

	return pp
}
