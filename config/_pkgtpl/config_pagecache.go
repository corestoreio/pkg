// +build ignore

package pagecache

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "system",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "full_page_cache",
				Label:     `Full Page Cache`,
				SortOrder: 600,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: system/full_page_cache/caching_application
						ID:      "caching_application",
						Label:   `Caching Application`,
						Type:    config.TypeSelect,
						Visible: config.VisibleYes,
						Scope:   scope.NewPerm(scope.DefaultID),
						Default: true,
						// SourceModel: Otnegam\PageCache\Model\System\Config\Source\Application
					},

					&config.Field{
						// Path: system/full_page_cache/ttl
						ID:        "ttl",
						Label:     `TTL for public content`,
						Comment:   element.LongText(`Public content cache lifetime in seconds. If field is empty default value 86400 will be saved.`),
						Type:      config.TypeText,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   86400,
						// BackendModel: Otnegam\PageCache\Model\System\Config\Backend\Ttl
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "system",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "full_page_cache",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: system/full_page_cache/varnish3
						ID:      `varnish3`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"path":"varnish3.vcl"}`,
					},

					&config.Field{
						// Path: system/full_page_cache/varnish4
						ID:      `varnish4`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"path":"varnish4.vcl"}`,
					},

					&config.Field{
						// Path: system/full_page_cache/default
						ID:      `default`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"access_list":"localhost","backend_host":"localhost","backend_port":"8080","ttl":"86400"}`,
					},
				),
			},
		),
	},
)
