// +build ignore

package cms

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "web",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "default",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: web/default/cms_home_page
						ID:        "cms_home_page",
						Label:     `CMS Home Page`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `home`,
						// SourceModel: Otnegam\Cms\Model\Config\Source\Page
					},

					&config.Field{
						// Path: web/default/cms_no_route
						ID:        "cms_no_route",
						Label:     `CMS No Route Page`,
						Type:      config.TypeSelect,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `no-route`,
						// SourceModel: Otnegam\Cms\Model\Config\Source\Page
					},

					&config.Field{
						// Path: web/default/cms_no_cookies
						ID:        "cms_no_cookies",
						Label:     `CMS No Cookies Page`,
						Type:      config.TypeSelect,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `enable-cookies`,
						// SourceModel: Otnegam\Cms\Model\Config\Source\Page
					},

					&config.Field{
						// Path: web/default/show_cms_breadcrumbs
						ID:        "show_cms_breadcrumbs",
						Label:     `Show Breadcrumbs for CMS Pages`,
						Type:      config.TypeSelect,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},

			&config.Group{
				ID:        "browser_capabilities",
				Label:     `Browser Capabilities Detection`,
				SortOrder: 200,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: web/browser_capabilities/cookies
						ID:        "cookies",
						Label:     `Redirect to CMS-page if Cookies are Disabled`,
						Type:      config.TypeSelect,
						SortOrder: 100,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: web/browser_capabilities/javascript
						ID:        "javascript",
						Label:     `Show Notice if JavaScript is Disabled`,
						Type:      config.TypeSelect,
						SortOrder: 200,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: web/browser_capabilities/local_storage
						ID:        "local_storage",
						Label:     `Show Notice if Local Storage is Disabled`,
						Type:      config.TypeSelect,
						SortOrder: 300,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},
	&config.Section{
		ID:        "cms",
		Label:     `Content Management`,
		SortOrder: 1001,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Cms::config_cms
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "wysiwyg",
				Label:     `WYSIWYG Options`,
				SortOrder: 100,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: cms/wysiwyg/enabled
						ID:        "enabled",
						Label:     `Enable WYSIWYG Editor`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `enabled`,
						// SourceModel: Otnegam\Cms\Model\Config\Source\Wysiwyg\Enabled
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "web",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "default",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: web/default/front
						ID:      `front`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `cms`,
					},

					&config.Field{
						// Path: web/default/no_route
						ID:      `no_route`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `cms/noroute/index`,
					},
				),
			},
		),
	},
	&config.Section{
		ID: "system",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "media_storage_configuration",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: system/media_storage_configuration/allowed_resources
						ID:      `allowed_resources`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"wysiwyg_image_folder":"wysiwyg"}`,
					},
				),
			},
		),
	},
)
