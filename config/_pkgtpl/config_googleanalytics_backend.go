// +build ignore

package googleanalytics

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
	// GoogleAnalyticsActive => Enable.
	// Path: google/analytics/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	GoogleAnalyticsActive model.Bool

	// GoogleAnalyticsAccount => Account Number.
	// Path: google/analytics/account
	GoogleAnalyticsAccount model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.GoogleAnalyticsActive = model.NewBool(`google/analytics/active`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.GoogleAnalyticsAccount = model.NewStr(`google/analytics/account`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
