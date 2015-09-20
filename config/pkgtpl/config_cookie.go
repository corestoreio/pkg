// +build ignore

package cookie

import "github.com/corestoreio/csfw/config"

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "web",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "cookie",
				Label:     `Default Cookie Settings`,
				Comment:   ``,
				SortOrder: 50,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `web/cookie/cookie_lifetime`,
						ID:           "cookie_lifetime",
						Label:        `Cookie Lifetime`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      3600,
						BackendModel: nil, // Magento\Cookie\Model\Config\Backend\Lifetime
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `web/cookie/cookie_path`,
						ID:           "cookie_path",
						Label:        `Cookie Path`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Cookie\Model\Config\Backend\Path
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `web/cookie/cookie_domain`,
						ID:           "cookie_domain",
						Label:        `Cookie Domain`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Cookie\Model\Config\Backend\Domain
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `web/cookie/cookie_httponly`,
						ID:           "cookie_httponly",
						Label:        `Use HTTP Only`,
						Comment:      `<strong style="color:red">Warning</strong>:  Do not set to "No". User security could be compromised.`,
						Type:         config.TypeSelect,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `web/cookie/cookie_restriction`,
						ID:           "cookie_restriction",
						Label:        `Cookie Restriction Mode`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      false,
						BackendModel: nil, // Magento\Cookie\Model\Config\Backend\Cookie
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "web",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "cookie",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `web/cookie/cookie_restriction_lifetime`,
						ID:      "cookie_restriction_lifetime",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: 31536000,
					},
				},
			},
		},
	},
)
