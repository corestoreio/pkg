// +build ignore

package pagecache

import "github.com/corestoreio/csfw/config"

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "system",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "full_page_cache",
				Label:     `Full Page Cache`,
				Comment:   ``,
				SortOrder: 600,
				Scope:     config.NewScopePerm(config.IDScopeDefault),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `system/full_page_cache/caching_application`,
						ID:           "caching_application",
						Label:        `Caching Application`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    0,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\PageCache\Model\System\Config\Source\Application
					},

					&config.Field{
						// Path: `system/full_page_cache/ttl`,
						ID:           "ttl",
						Label:        `TTL for public content`,
						Comment:      `Public content cache lifetime in seconds. If field is empty default value 120 will be saved.`,
						Type:         config.TypeText,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault),
						Default:      120,
						BackendModel: nil, // Magento\PageCache\Model\System\Config\Backend\Ttl
						SourceModel:  nil,
					},
				},
			},
		},
	},

	// Hidden Configuration
	&config.Section{
		ID: "system",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "full_page_cache",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `system/full_page_cache/varnish3`,
						ID:      "varnish3",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `{"path":"Magento\/PageCache\/etc\/varnish3.vcl"}`,
					},

					&config.Field{
						// Path: `system/full_page_cache/varnish4`,
						ID:      "varnish4",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `{"path":"Magento\/PageCache\/etc\/varnish4.vcl"}`,
					},

					&config.Field{
						// Path: `system/full_page_cache/default`,
						ID:      "default",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `{"access_list":"localhost","backend_host":"localhost","backend_port":"8080","ttl":"120"}`,
					},
				},
			},
		},
	},
)
