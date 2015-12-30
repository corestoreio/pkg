// +build ignore

package ui

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = element.MustNewConfiguration(
	&element.Section{
		ID: "dev",
		Groups: element.NewGroupSlice(
			&element.Group{
				ID: "js",
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: dev/js/session_storage_logging
						ID:        "session_storage_logging",
						Label:     `Log JS Errors to Session Storage`,
						Comment:   element.LongText(`If enabled, can be used by functional tests for extended reporting`),
						Type:      element.TypeSelect,
						SortOrder: 100,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&element.Field{
						// Path: dev/js/session_storage_key
						ID:        "session_storage_key",
						Label:     `Log JS Errors to Session Storage Key`,
						Comment:   element.LongText(`Use this key to retrieve collected js errors`),
						Type:      element.TypeText,
						SortOrder: 110,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   `collected_errors`,
					},
				),
			},
		),
	},
)
