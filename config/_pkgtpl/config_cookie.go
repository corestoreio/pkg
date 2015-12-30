// +build ignore

package cookie

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
				ID:        "cookie",
				Label:     `Default Cookie Settings`,
				SortOrder: 50,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: web/cookie/cookie_lifetime
						ID:        "cookie_lifetime",
						Label:     `Cookie Lifetime`,
						Type:      config.TypeText,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   3600,
						// BackendModel: Otnegam\Cookie\Model\Config\Backend\Lifetime
					},

					&config.Field{
						// Path: web/cookie/cookie_path
						ID:        "cookie_path",
						Label:     `Cookie Path`,
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Cookie\Model\Config\Backend\Path
					},

					&config.Field{
						// Path: web/cookie/cookie_domain
						ID:        "cookie_domain",
						Label:     `Cookie Domain`,
						Type:      config.TypeText,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Cookie\Model\Config\Backend\Domain
					},

					&config.Field{
						// Path: web/cookie/cookie_httponly
						ID:        "cookie_httponly",
						Label:     `Use HTTP Only`,
						Comment:   element.LongText(`<strong style="color:red">Warning</strong>: Do not set to "No". User security could be compromised.`),
						Type:      config.TypeSelect,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: web/cookie/cookie_restriction
						ID:        "cookie_restriction",
						Label:     `Cookie Restriction Mode`,
						Type:      config.TypeSelect,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// BackendModel: Otnegam\Cookie\Model\Config\Backend\Cookie
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
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
				ID: "cookie",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: web/cookie/cookie_restriction_lifetime
						ID:      `cookie_restriction_lifetime`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: 31536000,
					},
				),
			},
		),
	},
)
