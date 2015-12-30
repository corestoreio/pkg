// +build ignore

package vault

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
			ID: "payment",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "vault",
					Label:     `Vault Provider`,
					SortOrder: 2,
					Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: payment/vault/vault_payment
							ID:      "vault_payment",
							Label:   `Vault Provider`,
							Comment: element.LongText(`Specified provider should be enabled.`),
							Type:    element.TypeSelect,
							Visible: element.VisibleYes,
							Scope:   scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Otnegam\Vault\Model\Adminhtml\Source\VaultProvidersMap
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "payment",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "vault",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: payment/vault/debug
							ID:      `debug`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						&element.Field{
							// Path: payment/vault/model
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Otnegam\Vault\Model\VaultPaymentInterface`,
						},
					),
				},
			),
		},
	)
}
