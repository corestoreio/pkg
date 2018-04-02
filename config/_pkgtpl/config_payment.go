// +build ignore

package payment

import (
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.Sections

func init() {
	ConfigStructure = element.MustMakeSectionsValidate(
		element.Section{
			ID:        "payment",
			Label:     `Payment Methods`,
			SortOrder: 400,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Payment::payment
			Groups:    element.MakeGroups(),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: "payment",
			Groups: element.MakeGroups(
				element.Group{
					ID: "free",
					Fields: element.MakeFields(
						element.Field{
							// Path: payment/free/active
							ID:      `active`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/free/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\Payment\Model\Method\Free`,
						},

						element.Field{
							// Path: payment/free/order_status
							ID:      `order_status`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `pending`,
						},

						element.Field{
							// Path: payment/free/title
							ID:      `title`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `No Payment Information Required`,
						},

						element.Field{
							// Path: payment/free/allowspecific
							ID:      `allowspecific`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						element.Field{
							// Path: payment/free/sort_order
							ID:      `sort_order`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},
					),
				},

				element.Group{
					ID: "substitution",
					Fields: element.MakeFields(
						element.Field{
							// Path: payment/substitution/active
							ID:      `active`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						element.Field{
							// Path: payment/substitution/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\Payment\Model\Method\Substitution`,
						},

						element.Field{
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
	Backend = NewBackend(ConfigStructure)
}
