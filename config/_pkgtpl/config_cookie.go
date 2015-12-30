// +build ignore

package cookie

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package.
// Used in frontend and backend. See init() for details.
var PackageConfiguration element.SectionSlice

func init() {
	PackageConfiguration = element.MustNewConfiguration(
		&element.Section{
			ID: "web",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "cookie",
					Label:     `Default Cookie Settings`,
					SortOrder: 50,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: web/cookie/cookie_lifetime
							ID:        "cookie_lifetime",
							Label:     `Cookie Lifetime`,
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   3600,
							// BackendModel: Otnegam\Cookie\Model\Config\Backend\Lifetime
						},

						&element.Field{
							// Path: web/cookie/cookie_path
							ID:        "cookie_path",
							Label:     `Cookie Path`,
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// BackendModel: Otnegam\Cookie\Model\Config\Backend\Path
						},

						&element.Field{
							// Path: web/cookie/cookie_domain
							ID:        "cookie_domain",
							Label:     `Cookie Domain`,
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// BackendModel: Otnegam\Cookie\Model\Config\Backend\Domain
						},

						&element.Field{
							// Path: web/cookie/cookie_httponly
							ID:        "cookie_httponly",
							Label:     `Use HTTP Only`,
							Comment:   element.LongText(`<strong style="color:red">Warning</strong>: Do not set to "No". User security could be compromised.`),
							Type:      element.TypeSelect,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/cookie/cookie_restriction
							ID:        "cookie_restriction",
							Label:     `Cookie Restriction Mode`,
							Type:      element.TypeSelect,
							SortOrder: 50,
							Visible:   element.VisibleYes,
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
		&element.Section{
			ID: "web",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "cookie",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: web/cookie/cookie_restriction_lifetime
							ID:      `cookie_restriction_lifetime`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: 31536000,
						},
					),
				},
			),
		},
	)
}
