// +build ignore

package webapi

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "webapi",
		Label:     "Magento Web API",
		SortOrder: 102,
		Scope:     scope.PermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "soap",
				Label:     `SOAP Settings`,
				Comment:   ``,
				SortOrder: 1,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `webapi/soap/charset`,
						ID:           "charset",
						Label:        `Default Response Charset`,
						Comment:      `If empty, UTF-8 will be used.`,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},
				},
			},
		},
	},
)
