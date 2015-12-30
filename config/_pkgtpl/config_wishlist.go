// +build ignore

package wishlist

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "wishlist",
		Label:     `Wish List`,
		SortOrder: 140,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Wishlist::config_wishlist
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "email",
				Label:     `Share Options`,
				SortOrder: 2,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: wishlist/email/email_identity
						ID:        "email_identity",
						Label:     `Email Sender`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `general`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: wishlist/email/email_template
						ID:        "email_template",
						Label:     `Email Template`,
						Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `wishlist_email_email_template`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: wishlist/email/number_limit
						ID:        "number_limit",
						Label:     `Max Emails Allowed to be Sent`,
						Comment:   element.LongText(`10 by default. Max - 10000`),
						Type:      config.TypeText,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   10,
					},

					&config.Field{
						// Path: wishlist/email/text_limit
						ID:        "text_limit",
						Label:     `Email Text Length Limit`,
						Comment:   element.LongText(`255 by default`),
						Type:      config.TypeText,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   255,
					},
				),
			},

			&config.Group{
				ID:        "general",
				Label:     `General Options`,
				SortOrder: 1,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: wishlist/general/active
						ID:        "active",
						Label:     `Enabled`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},

			&config.Group{
				ID:        "wishlist_link",
				Label:     `My Wish List Link`,
				SortOrder: 3,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: wishlist/wishlist_link/use_qty
						ID:        "use_qty",
						Label:     `Display Wish List Summary`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Wishlist\Model\Config\Source\Summary
					},
				),
			},
		),
	},
	&config.Section{
		ID: "rss",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "wishlist",
				Label:     `Wish List`,
				SortOrder: 2,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: rss/wishlist/active
						ID:        "active",
						Label:     `Enable RSS`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
					},
				),
			},
		),
	},
)
