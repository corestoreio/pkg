// +build ignore

package webapi

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
			ID:        "webapi",
			Label:     `Otnegam Web API`,
			SortOrder: 102,
			Scope:     scope.PermAll,
			Resource:  0, // Otnegam_Webapi::config_webapi
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "soap",
					Label:     `SOAP Settings`,
					SortOrder: 1,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: webapi/soap/charset
							ID:        "charset",
							Label:     `Default Response Charset`,
							Comment:   element.LongText(`If empty, UTF-8 will be used.`),
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
