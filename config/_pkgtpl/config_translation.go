// +build ignore

package translation

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "dev",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "js",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: dev/js/translate_strategy
						ID:        "translate_strategy",
						Label:     `Translation Strategy`,
						Comment:   element.LongText(`Please put your store into maintenance mode and redeploy static files after changing strategy`),
						Type:      config.TypeSelect,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   `dictionary`,
						// SourceModel: Otnegam\Translation\Model\Js\Config\Source\Strategy
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "dev",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "translate_inline",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: dev/translate_inline/active
						ID:      `active`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: dev/translate_inline/active_admin
						ID:      `active_admin`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: dev/translate_inline/invalid_caches
						ID:      `invalid_caches`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"block_html":null}`,
					},
				),
			},
		),
	},
)
