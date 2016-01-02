// +build ignore

package payment

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
			ID:        "payment",
			Label:     `Payment Methods`,
			SortOrder: 400,
			Scope:     scope.PermAll,
			Resource:  0, // Otnegam_Payment::payment
			Groups:    element.NewGroupSlice(),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "payment",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "free",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: payment/free/active
							ID:      `active`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						&element.Field{
							// Path: payment/free/model
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Otnegam\Payment\Model\Method\Free`,
						},

						&element.Field{
							// Path: payment/free/order_status
							ID:      `order_status`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `pending`,
						},

						&element.Field{
							// Path: payment/free/title
							ID:      `title`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `No Payment Information Required`,
						},

						&element.Field{
							// Path: payment/free/allowspecific
							ID:      `allowspecific`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						&element.Field{
							// Path: payment/free/sort_order
							ID:      `sort_order`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},
					),
				},

				&element.Group{
					ID: "substitution",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: payment/substitution/active
							ID:      `active`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						&element.Field{
							// Path: payment/substitution/model
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Otnegam\Payment\Model\Method\Substitution`,
						},

						&element.Field{
							// Path: payment/substitution/allowspecific
							ID:      `allowspecific`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},
					),
				},
			),
		},
	)
	Path = NewPath(PackageConfiguration)
}
