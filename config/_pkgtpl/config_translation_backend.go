// +build ignore

package translation

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
	// DevJsTranslateStrategy => Translation Strategy.
	// Please put your store into maintenance mode and redeploy static files after
	// changing strategy
	// Path: dev/js/translate_strategy
	// SourceModel: Otnegam\Translation\Model\Js\Config\Source\Strategy
	DevJsTranslateStrategy model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.DevJsTranslateStrategy = model.NewStr(`dev/js/translate_strategy`, model.WithConfigStructure(cfgStruct))

	return pp
}
