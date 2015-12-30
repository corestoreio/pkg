// +build ignore

package review

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "catalog",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "review",
				Label:     `Product Reviews`,
				SortOrder: 100,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/review/allow_guest
						ID:        "allow_guest",
						Label:     `Allow Guests to Write Reviews`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},
)
