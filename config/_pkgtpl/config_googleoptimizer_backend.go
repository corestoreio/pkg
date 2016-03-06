// +build ignore

package googleoptimizer

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
	// GoogleAnalyticsExperiments => Enable Content Experiments.
	// Path: google/analytics/experiments
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	GoogleAnalyticsExperiments model.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.GoogleAnalyticsExperiments = model.NewBool(`google/analytics/experiments`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
