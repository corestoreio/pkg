// +build ignore

package googleanalytics

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
	// GoogleAnalyticsActive => Enable.
	// Path: google/analytics/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	GoogleAnalyticsActive cfgmodel.Bool

	// GoogleAnalyticsAccount => Account Number.
	// Path: google/analytics/account
	GoogleAnalyticsAccount cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.GoogleAnalyticsActive = cfgmodel.NewBool(`google/analytics/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.GoogleAnalyticsAccount = cfgmodel.NewStr(`google/analytics/account`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
