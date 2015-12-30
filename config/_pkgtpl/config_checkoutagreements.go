// +build ignore

package checkoutagreements

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "checkout",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "options",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: checkout/options/enable_agreements
						ID:        "enable_agreements",
						Label:     `Enable Terms and Conditions`,
						Type:      config.TypeSelect,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},
)
