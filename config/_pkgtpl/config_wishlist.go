// +build ignore

package wishlist

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
			ID:        "wishlist",
			Label:     `Wish List`,
			SortOrder: 140,
			Scope:     scope.PermAll,
			Resource:  0, // Magento_Wishlist::config_wishlist
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "email",
					Label:     `Share Options`,
					SortOrder: 2,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: wishlist/email/email_identity
							ID:        "email_identity",
							Label:     `Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `general`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},

						&element.Field{
							// Path: wishlist/email/email_template
							ID:        "email_template",
							Label:     `Email Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `wishlist_email_email_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: wishlist/email/number_limit
							ID:        "number_limit",
							Label:     `Max Emails Allowed to be Sent`,
							Comment:   text.Long(`10 by default. Max - 10000`),
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   10,
						},

						&element.Field{
							// Path: wishlist/email/text_limit
							ID:        "text_limit",
							Label:     `Email Text Length Limit`,
							Comment:   text.Long(`255 by default`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   255,
						},
					),
				},

				&element.Group{
					ID:        "general",
					Label:     `General Options`,
					SortOrder: 1,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: wishlist/general/active
							ID:        "active",
							Label:     `Enabled`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        "wishlist_link",
					Label:     `My Wish List Link`,
					SortOrder: 3,
					Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: wishlist/wishlist_link/use_qty
							ID:        "use_qty",
							Label:     `Display Wish List Summary`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Magento\Wishlist\Model\Config\Source\Summary
						},
					),
				},
			),
		},
		&element.Section{
			ID: "rss",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "wishlist",
					Label:     `Wish List`,
					SortOrder: 2,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: rss/wishlist/active
							ID:        "active",
							Label:     `Enable RSS`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
