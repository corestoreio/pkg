// +build ignore

package vault

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "payment",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "vault",
				Label:     `Vault Provider`,
				SortOrder: 2,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/vault/vault_payment
						ID:      "vault_payment",
						Label:   `Vault Provider`,
						Comment: element.LongText(`Specified provider should be enabled.`),
						Type:    config.TypeSelect,
						Visible: config.VisibleYes,
						Scope:   scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Vault\Model\Adminhtml\Source\VaultProvidersMap
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "payment",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "vault",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/vault/debug
						ID:      `debug`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/vault/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\Vault\Model\VaultPaymentInterface`,
					},
				),
			},
		),
	},
)
