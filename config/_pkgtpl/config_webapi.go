// +build ignore

package webapi

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "webapi",
		Label:     `Otnegam Web API`,
		SortOrder: 102,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Webapi::config_webapi
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "soap",
				Label:     `SOAP Settings`,
				SortOrder: 1,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: webapi/soap/charset
						ID:        "charset",
						Label:     `Default Response Charset`,
						Comment:   element.LongText(`If empty, UTF-8 will be used.`),
						Type:      config.TypeText,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},
				),
			},
		),
	},
)
