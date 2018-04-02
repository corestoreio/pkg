// +build ignore

package translation

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
	// DevJsTranslateStrategy => Translation Strategy.
	// Please put your store into maintenance mode and redeploy static files after
	// changing strategy
	// Path: dev/js/translate_strategy
	// SourceModel: Magento\Translation\Model\Js\Config\Source\Strategy
	DevJsTranslateStrategy cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.DevJsTranslateStrategy = cfgmodel.NewStr(`dev/js/translate_strategy`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
