// +build ignore

package vault

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
			ID: "payment",
			Groups: element.MakeGroups(
				element.Group{
					ID:        "vault",
					Label:     `Vault Provider`,
					SortOrder: 2,
					Scopes:    scope.PermWebsite,
					Fields: element.MakeFields(
						element.Field{
							// Path: payment/vault/vault_payment
							ID:      "vault_payment",
							Label:   `Vault Provider`,
							Comment: text.Long(`Specified provider should be enabled.`),
							Type:    element.TypeSelect,
							Visible: element.VisibleYes,
							Scopes:  scope.PermWebsite,
							// SourceModel: Magento\Vault\Model\Adminhtml\Source\VaultProvidersMap
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: "payment",
			Groups: element.MakeGroups(
				element.Group{
					ID: "vault",
					Fields: element.MakeFields(
						element.Field{
							// Path: payment/vault/debug
							ID:      `debug`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/vault/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\Vault\Model\VaultPaymentInterface`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
