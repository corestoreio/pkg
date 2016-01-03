// +build ignore

package cms

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
			ID: "web",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "default",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: web/default/cms_home_page
							ID:        "cms_home_page",
							Label:     `CMS Home Page`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `home`,
							// SourceModel: Otnegam\Cms\Model\Config\Source\Page
						},

						&element.Field{
							// Path: web/default/cms_no_route
							ID:        "cms_no_route",
							Label:     `CMS No Route Page`,
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `no-route`,
							// SourceModel: Otnegam\Cms\Model\Config\Source\Page
						},

						&element.Field{
							// Path: web/default/cms_no_cookies
							ID:        "cms_no_cookies",
							Label:     `CMS No Cookies Page`,
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `enable-cookies`,
							// SourceModel: Otnegam\Cms\Model\Config\Source\Page
						},

						&element.Field{
							// Path: web/default/show_cms_breadcrumbs
							ID:        "show_cms_breadcrumbs",
							Label:     `Show Breadcrumbs for CMS Pages`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        "browser_capabilities",
					Label:     `Browser Capabilities Detection`,
					SortOrder: 200,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: web/browser_capabilities/cookies
							ID:        "cookies",
							Label:     `Redirect to CMS-page if Cookies are Disabled`,
							Type:      element.TypeSelect,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/browser_capabilities/javascript
							ID:        "javascript",
							Label:     `Show Notice if JavaScript is Disabled`,
							Type:      element.TypeSelect,
							SortOrder: 200,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/browser_capabilities/local_storage
							ID:        "local_storage",
							Label:     `Show Notice if Local Storage is Disabled`,
							Type:      element.TypeSelect,
							SortOrder: 300,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
		&element.Section{
			ID:        "cms",
			Label:     `Content Management`,
			SortOrder: 1001,
			Scope:     scope.PermAll,
			Resource:  0, // Otnegam_Cms::config_cms
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "wysiwyg",
					Label:     `WYSIWYG Options`,
					SortOrder: 100,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: cms/wysiwyg/enabled
							ID:        "enabled",
							Label:     `Enable WYSIWYG Editor`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `enabled`,
							// SourceModel: Otnegam\Cms\Model\Config\Source\Wysiwyg\Enabled
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "web",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "default",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: web/default/front
							ID:      `front`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `cms`,
						},

						&element.Field{
							// Path: web/default/no_route
							ID:      `no_route`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `cms/noroute/index`,
						},
					),
				},
			),
		},
		&element.Section{
			ID: "system",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "media_storage_configuration",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: system/media_storage_configuration/allowed_resources
							ID:      `allowed_resources`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"wysiwyg_image_folder":"wysiwyg"}`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
