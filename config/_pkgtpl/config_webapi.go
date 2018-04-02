// +build ignore

package webapi

import (
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.Sections

func init() {
	ConfigStructure = element.MustMakeSectionsValidate(
		element.Section{
			ID:        "webapi",
			Label:     `Magento Web API`,
			SortOrder: 102,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Webapi::config_webapi
			Groups: element.MakeGroups(
				element.Group{
					ID:        "soap",
					Label:     `SOAP Settings`,
					SortOrder: 1,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: webapi/soap/charset
							ID:        "charset",
							Label:     `Default Response Charset`,
							Comment:   text.Long(`If empty, UTF-8 will be used.`),
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
