// +build ignore

package ui

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathDevJsSessionStorageLogging => Log JS Errors to Session Storage.
// If enabled, can be used by functional tests for extended reporting
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathDevJsSessionStorageLogging = model.NewBool(`dev/js/session_storage_logging`, model.WithPkgCfg(PackageConfiguration))

// PathDevJsSessionStorageKey => Log JS Errors to Session Storage Key.
// Use this key to retrieve collected js errors
var PathDevJsSessionStorageKey = model.NewStr(`dev/js/session_storage_key`, model.WithPkgCfg(PackageConfiguration))
