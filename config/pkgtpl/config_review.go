// +build ignore

package review

import "github.com/corestoreio/csfw/config"

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "catalog",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "review",
				Label:     `Product Reviews`,
				Comment:   ``,
				SortOrder: 100,
				Scope:     config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/review/allow_guest`,
						ID:           "allow_guest",
						Label:        `Allow Guests to Write Reviews`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},
)
