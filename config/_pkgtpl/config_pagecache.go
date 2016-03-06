// +build ignore

package pagecache

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
			ID: "system",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "full_page_cache",
					Label:     `Full Page Cache`,
					SortOrder: 600,
					Scopes:    scope.PermDefault,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: system/full_page_cache/caching_application
							ID:      "caching_application",
							Label:   `Caching Application`,
							Type:    element.TypeSelect,
							Visible: element.VisibleYes,
							Scopes:  scope.PermDefault,
							Default: true,
							// SourceModel: Magento\PageCache\Model\System\Config\Source\Application
						},

						&element.Field{
							// Path: system/full_page_cache/ttl
							ID:        "ttl",
							Label:     `TTL for public content`,
							Comment:   text.Long(`Public content cache lifetime in seconds. If field is empty default value 86400 will be saved.`),
							Type:      element.TypeText,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   86400,
							// BackendModel: Magento\PageCache\Model\System\Config\Backend\Ttl
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "system",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "full_page_cache",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: system/full_page_cache/varnish3
							ID:      `varnish3`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"path":"varnish3.vcl"}`,
						},

						&element.Field{
							// Path: system/full_page_cache/varnish4
							ID:      `varnish4`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"path":"varnish4.vcl"}`,
						},

						&element.Field{
							// Path: system/full_page_cache/default
							ID:      `default`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"access_list":"localhost","backend_host":"localhost","backend_port":"8080","ttl":"86400"}`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
