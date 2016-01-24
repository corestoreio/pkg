// +build ignore

package ui

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
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
							Comment:   text.Long(`If enabled, can be used by functional tests for extended reporting`),
							Type:      element.TypeSelect,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: dev/js/session_storage_key
							ID:        "session_storage_key",
							Label:     `Log JS Errors to Session Storage Key`,
							Comment:   text.Long(`Use this key to retrieve collected js errors`),
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
	Backend = NewBackend(ConfigStructure)
}
