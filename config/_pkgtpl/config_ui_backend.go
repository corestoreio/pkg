// +build ignore

package ui

import (
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// DevJsSessionStorageLogging => Log JS Errors to Session Storage.
	// If enabled, can be used by functional tests for extended reporting
	// Path: dev/js/session_storage_logging
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	DevJsSessionStorageLogging cfgmodel.Bool

	// DevJsSessionStorageKey => Log JS Errors to Session Storage Key.
	// Use this key to retrieve collected js errors
	// Path: dev/js/session_storage_key
	DevJsSessionStorageKey cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.DevJsSessionStorageLogging = cfgmodel.NewBool(`dev/js/session_storage_logging`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DevJsSessionStorageKey = cfgmodel.NewStr(`dev/js/session_storage_key`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
