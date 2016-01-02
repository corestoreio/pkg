// +build ignore

package translation

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
	// DevJsTranslateStrategy => Translation Strategy.
	// Please put your store into maintenance mode and redeploy static files after
	// changing strategy
	// Path: dev/js/translate_strategy
	// SourceModel: Otnegam\Translation\Model\Js\Config\Source\Strategy
	DevJsTranslateStrategy model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.DevJsTranslateStrategy = model.NewStr(`dev/js/translate_strategy`, model.WithPkgCfg(pkgCfg))

	return pp
}
