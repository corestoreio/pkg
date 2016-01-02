// +build ignore

package googleanalytics

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
	// GoogleAnalyticsActive => Enable.
	// Path: google/analytics/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	GoogleAnalyticsActive model.Bool

	// GoogleAnalyticsAccount => Account Number.
	// Path: google/analytics/account
	GoogleAnalyticsAccount model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.GoogleAnalyticsActive = model.NewBool(`google/analytics/active`, model.WithPkgCfg(pkgCfg))
	pp.GoogleAnalyticsAccount = model.NewStr(`google/analytics/account`, model.WithPkgCfg(pkgCfg))

	return pp
}
