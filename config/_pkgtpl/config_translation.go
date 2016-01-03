// +build ignore

package translation

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
			ID: "dev",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "js",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: dev/js/translate_strategy
							ID:        "translate_strategy",
							Label:     `Translation Strategy`,
							Comment:   element.LongText(`Please put your store into maintenance mode and redeploy static files after changing strategy`),
							Type:      element.TypeSelect,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   `dictionary`,
							// SourceModel: Otnegam\Translation\Model\Js\Config\Source\Strategy
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "dev",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "translate_inline",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: dev/translate_inline/active
							ID:      `active`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						&element.Field{
							// Path: dev/translate_inline/active_admin
							ID:      `active_admin`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						&element.Field{
							// Path: dev/translate_inline/invalid_caches
							ID:      `invalid_caches`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"block_html":null}`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
