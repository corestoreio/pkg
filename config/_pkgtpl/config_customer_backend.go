// +build ignore

package customer

import (
	"github.com/corestoreio/pkg/config/cfgmodel"
	"github.com/corestoreio/pkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// CustomerAccountShareScope => Share Customer Accounts.
	// Path: customer/account_share/scope
	// BackendModel: Magento\Customer\Model\Config\Share
	// SourceModel: Magento\Customer\Model\Config\Share
	CustomerAccountShareScope cfgmodel.Str

	// CustomerCreateAccountAutoGroupAssign => Enable Automatic Assignment to Customer Group.
	// Path: customer/create_account/auto_group_assign
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CustomerCreateAccountAutoGroupAssign cfgmodel.Bool

	// CustomerCreateAccountTaxCalculationAddressType => Tax Calculation Based On.
	// Path: customer/create_account/tax_calculation_address_type
	// SourceModel: Magento\Customer\Model\Config\Source\Address\Type
	CustomerCreateAccountTaxCalculationAddressType cfgmodel.Str

	// CustomerCreateAccountDefaultGroup => Default Group.
	// Path: customer/create_account/default_group
	// SourceModel: Magento\Customer\Model\Config\Source\Group
	CustomerCreateAccountDefaultGroup cfgmodel.Str

	// CustomerCreateAccountVivDomesticGroup => Group for Valid VAT ID - Domestic.
	// Path: customer/create_account/viv_domestic_group
	// SourceModel: Magento\Customer\Model\Config\Source\Group
	CustomerCreateAccountVivDomesticGroup cfgmodel.Str

	// CustomerCreateAccountVivIntraUnionGroup => Group for Valid VAT ID - Intra-Union.
	// Path: customer/create_account/viv_intra_union_group
	// SourceModel: Magento\Customer\Model\Config\Source\Group
	CustomerCreateAccountVivIntraUnionGroup cfgmodel.Str

	// CustomerCreateAccountVivInvalidGroup => Group for Invalid VAT ID.
	// Path: customer/create_account/viv_invalid_group
	// SourceModel: Magento\Customer\Model\Config\Source\Group
	CustomerCreateAccountVivInvalidGroup cfgmodel.Str

	// CustomerCreateAccountVivErrorGroup => Validation Error Group.
	// Path: customer/create_account/viv_error_group
	// SourceModel: Magento\Customer\Model\Config\Source\Group
	CustomerCreateAccountVivErrorGroup cfgmodel.Str

	// CustomerCreateAccountVivOnEachTransaction => Validate on Each Transaction.
	// Path: customer/create_account/viv_on_each_transaction
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CustomerCreateAccountVivOnEachTransaction cfgmodel.Bool

	// CustomerCreateAccountVivDisableAutoGroupAssignDefault => Default Value for Disable Automatic Group Changes Based on VAT ID.
	// Path: customer/create_account/viv_disable_auto_group_assign_default
	// BackendModel: Magento\Customer\Model\Config\Backend\CreateAccount\DisableAutoGroupAssignDefault
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CustomerCreateAccountVivDisableAutoGroupAssignDefault cfgmodel.Bool

	// CustomerCreateAccountVatFrontendVisibility => Show VAT Number on Storefront.
	// To show VAT number on Storefront, set Show VAT Number on Storefront option
	// to Yes.
	// Path: customer/create_account/vat_frontend_visibility
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CustomerCreateAccountVatFrontendVisibility cfgmodel.Bool

	// CustomerCreateAccountEmailDomain => Default Email Domain.
	// Path: customer/create_account/email_domain
	CustomerCreateAccountEmailDomain cfgmodel.Str

	// CustomerCreateAccountEmailTemplate => Default Welcome Email.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: customer/create_account/email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	CustomerCreateAccountEmailTemplate cfgmodel.Str

	// CustomerCreateAccountEmailNoPasswordTemplate => Default Welcome Email Without Password.
	// This email will be sent instead of the Default Welcome Email, if a customer
	// was created without password.  Email template chosen based on theme
	// fallback when "Default" option is selected.
	// Path: customer/create_account/email_no_password_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	CustomerCreateAccountEmailNoPasswordTemplate cfgmodel.Str

	// CustomerCreateAccountEmailIdentity => Email Sender.
	// Path: customer/create_account/email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	CustomerCreateAccountEmailIdentity cfgmodel.Str

	// CustomerCreateAccountConfirm => Require Emails Confirmation.
	// Path: customer/create_account/confirm
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CustomerCreateAccountConfirm cfgmodel.Bool

	// CustomerCreateAccountEmailConfirmationTemplate => Confirmation Link Email.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: customer/create_account/email_confirmation_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	CustomerCreateAccountEmailConfirmationTemplate cfgmodel.Str

	// CustomerCreateAccountEmailConfirmedTemplate => Welcome Email.
	// This email will be sent instead of the Default Welcome Email, after account
	// confirmation.  Email template chosen based on theme fallback when "Default"
	// option is selected.
	// Path: customer/create_account/email_confirmed_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	CustomerCreateAccountEmailConfirmedTemplate cfgmodel.Str

	// CustomerCreateAccountGenerateHumanFriendlyId => Generate Human-Friendly Customer ID.
	// Path: customer/create_account/generate_human_friendly_id
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CustomerCreateAccountGenerateHumanFriendlyId cfgmodel.Bool

	// CustomerPasswordForgotEmailTemplate => Forgot Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: customer/password/forgot_email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	CustomerPasswordForgotEmailTemplate cfgmodel.Str

	// CustomerPasswordRemindEmailTemplate => Remind Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: customer/password/remind_email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	CustomerPasswordRemindEmailTemplate cfgmodel.Str

	// CustomerPasswordResetPasswordTemplate => Reset Password Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: customer/password/reset_password_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	CustomerPasswordResetPasswordTemplate cfgmodel.Str

	// CustomerPasswordForgotEmailIdentity => Password Template Email Sender.
	// Path: customer/password/forgot_email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	CustomerPasswordForgotEmailIdentity cfgmodel.Str

	// CustomerPasswordResetLinkExpirationPeriod => Recovery Link Expiration Period (days).
	// Please enter a number 1 or greater in this field.
	// Path: customer/password/reset_link_expiration_period
	// BackendModel: Magento\Customer\Model\Config\Backend\Password\Link\Expirationperiod
	CustomerPasswordResetLinkExpirationPeriod cfgmodel.Str

	// CustomerAddressStreetLines => Number of Lines in a Street Address.
	// Leave empty for default (2). Valid range: 1-4
	// Path: customer/address/street_lines
	// BackendModel: Magento\Customer\Model\Config\Backend\Address\Street
	CustomerAddressStreetLines cfgmodel.Str

	// CustomerAddressPrefixShow => Show Prefix.
	// The title that goes before name (Mr., Mrs., etc.)
	// Path: customer/address/prefix_show
	// BackendModel: Magento\Customer\Model\Config\Backend\Show\Address
	// SourceModel: Magento\Config\Model\Config\Source\Nooptreq
	CustomerAddressPrefixShow cfgmodel.Str

	// CustomerAddressPrefixOptions => Prefix Dropdown Options.
	// Semicolon (;) separated values.Put semicolon in the beginning for empty
	// first option.Leave empty for open text field.
	// Path: customer/address/prefix_options
	CustomerAddressPrefixOptions cfgmodel.Str

	// CustomerAddressMiddlenameShow => Show Middle Name (initial).
	// Always optional.
	// Path: customer/address/middlename_show
	// BackendModel: Magento\Customer\Model\Config\Backend\Show\Address
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CustomerAddressMiddlenameShow cfgmodel.Bool

	// CustomerAddressSuffixShow => Show Suffix.
	// The suffix that goes after name (Jr., Sr., etc.)
	// Path: customer/address/suffix_show
	// BackendModel: Magento\Customer\Model\Config\Backend\Show\Address
	// SourceModel: Magento\Config\Model\Config\Source\Nooptreq
	CustomerAddressSuffixShow cfgmodel.Str

	// CustomerAddressSuffixOptions => Suffix Dropdown Options.
	// Semicolon (;) separated values.Put semicolon in the beginning for empty
	// first option.Leave empty for open text field.
	// Path: customer/address/suffix_options
	CustomerAddressSuffixOptions cfgmodel.Str

	// CustomerAddressDobShow => Show Date of Birth.
	// Path: customer/address/dob_show
	// BackendModel: Magento\Customer\Model\Config\Backend\Show\Customer
	// SourceModel: Magento\Config\Model\Config\Source\Nooptreq
	CustomerAddressDobShow cfgmodel.Str

	// CustomerAddressTaxvatShow => Show Tax/VAT Number.
	// Path: customer/address/taxvat_show
	// BackendModel: Magento\Customer\Model\Config\Backend\Show\Customer
	// SourceModel: Magento\Config\Model\Config\Source\Nooptreq
	CustomerAddressTaxvatShow cfgmodel.Str

	// CustomerAddressGenderShow => Show Gender.
	// Path: customer/address/gender_show
	// BackendModel: Magento\Customer\Model\Config\Backend\Show\Customer
	// SourceModel: Magento\Config\Model\Config\Source\Nooptreq
	CustomerAddressGenderShow cfgmodel.Str

	// CustomerStartupRedirectDashboard => Redirect Customer to Account Dashboard after Logging in.
	// Customer will stay on the current page if "No" is selected.
	// Path: customer/startup/redirect_dashboard
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CustomerStartupRedirectDashboard cfgmodel.Bool

	// CustomerAddressTemplatesText => Text.
	// Path: customer/address_templates/text
	CustomerAddressTemplatesText cfgmodel.Str

	// CustomerAddressTemplatesOneline => Text One Line.
	// Path: customer/address_templates/oneline
	CustomerAddressTemplatesOneline cfgmodel.Str

	// CustomerAddressTemplatesHtml => HTML.
	// Path: customer/address_templates/html
	CustomerAddressTemplatesHtml cfgmodel.Str

	// CustomerAddressTemplatesPdf => PDF.
	// Path: customer/address_templates/pdf
	CustomerAddressTemplatesPdf cfgmodel.Str

	// CustomerOnlineCustomersOnlineMinutesInterval => Online Minutes Interval.
	// Leave empty for default (15 minutes).
	// Path: customer/online_customers/online_minutes_interval
	CustomerOnlineCustomersOnlineMinutesInterval cfgmodel.Str

	// GeneralStoreInformationValidateVatNumber => .
	// Path: general/store_information/validate_vat_number
	GeneralStoreInformationValidateVatNumber cfgmodel.Str

	// GeneralRestrictionAutocompleteOnStorefront => Enable Autocomplete on login/forgot password forms.
	// Path: general/restriction/autocomplete_on_storefront
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	GeneralRestrictionAutocompleteOnStorefront cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CustomerAccountShareScope = cfgmodel.NewStr(`customer/account_share/scope`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountAutoGroupAssign = cfgmodel.NewBool(`customer/create_account/auto_group_assign`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountTaxCalculationAddressType = cfgmodel.NewStr(`customer/create_account/tax_calculation_address_type`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountDefaultGroup = cfgmodel.NewStr(`customer/create_account/default_group`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountVivDomesticGroup = cfgmodel.NewStr(`customer/create_account/viv_domestic_group`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountVivIntraUnionGroup = cfgmodel.NewStr(`customer/create_account/viv_intra_union_group`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountVivInvalidGroup = cfgmodel.NewStr(`customer/create_account/viv_invalid_group`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountVivErrorGroup = cfgmodel.NewStr(`customer/create_account/viv_error_group`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountVivOnEachTransaction = cfgmodel.NewBool(`customer/create_account/viv_on_each_transaction`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountVivDisableAutoGroupAssignDefault = cfgmodel.NewBool(`customer/create_account/viv_disable_auto_group_assign_default`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountVatFrontendVisibility = cfgmodel.NewBool(`customer/create_account/vat_frontend_visibility`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountEmailDomain = cfgmodel.NewStr(`customer/create_account/email_domain`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountEmailTemplate = cfgmodel.NewStr(`customer/create_account/email_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountEmailNoPasswordTemplate = cfgmodel.NewStr(`customer/create_account/email_no_password_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountEmailIdentity = cfgmodel.NewStr(`customer/create_account/email_identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountConfirm = cfgmodel.NewBool(`customer/create_account/confirm`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountEmailConfirmationTemplate = cfgmodel.NewStr(`customer/create_account/email_confirmation_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountEmailConfirmedTemplate = cfgmodel.NewStr(`customer/create_account/email_confirmed_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCreateAccountGenerateHumanFriendlyId = cfgmodel.NewBool(`customer/create_account/generate_human_friendly_id`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerPasswordForgotEmailTemplate = cfgmodel.NewStr(`customer/password/forgot_email_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerPasswordRemindEmailTemplate = cfgmodel.NewStr(`customer/password/remind_email_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerPasswordResetPasswordTemplate = cfgmodel.NewStr(`customer/password/reset_password_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerPasswordForgotEmailIdentity = cfgmodel.NewStr(`customer/password/forgot_email_identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerPasswordResetLinkExpirationPeriod = cfgmodel.NewStr(`customer/password/reset_link_expiration_period`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerAddressStreetLines = cfgmodel.NewStr(`customer/address/street_lines`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerAddressPrefixShow = cfgmodel.NewStr(`customer/address/prefix_show`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerAddressPrefixOptions = cfgmodel.NewStr(`customer/address/prefix_options`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerAddressMiddlenameShow = cfgmodel.NewBool(`customer/address/middlename_show`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerAddressSuffixShow = cfgmodel.NewStr(`customer/address/suffix_show`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerAddressSuffixOptions = cfgmodel.NewStr(`customer/address/suffix_options`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerAddressDobShow = cfgmodel.NewStr(`customer/address/dob_show`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerAddressTaxvatShow = cfgmodel.NewStr(`customer/address/taxvat_show`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerAddressGenderShow = cfgmodel.NewStr(`customer/address/gender_show`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerStartupRedirectDashboard = cfgmodel.NewBool(`customer/startup/redirect_dashboard`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerAddressTemplatesText = cfgmodel.NewStr(`customer/address_templates/text`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerAddressTemplatesOneline = cfgmodel.NewStr(`customer/address_templates/oneline`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerAddressTemplatesHtml = cfgmodel.NewStr(`customer/address_templates/html`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerAddressTemplatesPdf = cfgmodel.NewStr(`customer/address_templates/pdf`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerOnlineCustomersOnlineMinutesInterval = cfgmodel.NewStr(`customer/online_customers/online_minutes_interval`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.GeneralStoreInformationValidateVatNumber = cfgmodel.NewStr(`general/store_information/validate_vat_number`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.GeneralRestrictionAutocompleteOnStorefront = cfgmodel.NewBool(`general/restriction/autocomplete_on_storefront`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
