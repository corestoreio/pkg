// +build ignore

package checkoutagreements

import "github.com/corestoreio/csfw/config"

import "github.com/corestoreio/csfw/config"

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "checkout",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "options",
				Label:     ``,
				Comment:   ``,
				SortOrder: 0,
				Scope:     config.NewScopePerm(),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `checkout/options/enable_agreements`,
						ID:           "enable_agreements",
						Label:        `Enable Terms and Conditions`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},
)
