// +build ignore

package ui

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
	// DevJsSessionStorageLogging => Log JS Errors to Session Storage.
	// If enabled, can be used by functional tests for extended reporting
	// Path: dev/js/session_storage_logging
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	DevJsSessionStorageLogging model.Bool

	// DevJsSessionStorageKey => Log JS Errors to Session Storage Key.
	// Use this key to retrieve collected js errors
	// Path: dev/js/session_storage_key
	DevJsSessionStorageKey model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.DevJsSessionStorageLogging = model.NewBool(`dev/js/session_storage_logging`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.DevJsSessionStorageKey = model.NewStr(`dev/js/session_storage_key`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
