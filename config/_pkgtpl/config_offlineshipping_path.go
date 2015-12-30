// +build ignore

package offlineshipping

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCarriersFlatrateActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersFlatrateActive = model.NewBool(`carriers/flatrate/active`)

// PathCarriersFlatrateName => Method Name.
var PathCarriersFlatrateName = model.NewStr(`carriers/flatrate/name`)

// PathCarriersFlatratePrice => Price.
var PathCarriersFlatratePrice = model.NewStr(`carriers/flatrate/price`)

// PathCarriersFlatrateHandlingType => Calculate Handling Fee.
// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
var PathCarriersFlatrateHandlingType = model.NewStr(`carriers/flatrate/handling_type`)

// PathCarriersFlatrateHandlingFee => Handling Fee.
var PathCarriersFlatrateHandlingFee = model.NewStr(`carriers/flatrate/handling_fee`)

// PathCarriersFlatrateSortOrder => Sort Order.
var PathCarriersFlatrateSortOrder = model.NewStr(`carriers/flatrate/sort_order`)

// PathCarriersFlatrateTitle => Title.
var PathCarriersFlatrateTitle = model.NewStr(`carriers/flatrate/title`)

// PathCarriersFlatrateType => Type.
// SourceModel: Otnegam\OfflineShipping\Model\Config\Source\Flatrate
var PathCarriersFlatrateType = model.NewStr(`carriers/flatrate/type`)

// PathCarriersFlatrateSallowspecific => Ship to Applicable Countries.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
var PathCarriersFlatrateSallowspecific = model.NewStr(`carriers/flatrate/sallowspecific`)

// PathCarriersFlatrateSpecificcountry => Ship to Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathCarriersFlatrateSpecificcountry = model.NewStringCSV(`carriers/flatrate/specificcountry`)

// PathCarriersFlatrateShowmethod => Show Method if Not Applicable.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersFlatrateShowmethod = model.NewBool(`carriers/flatrate/showmethod`)

// PathCarriersFlatrateSpecificerrmsg => Displayed Error Message.
var PathCarriersFlatrateSpecificerrmsg = model.NewStr(`carriers/flatrate/specificerrmsg`)

// PathCarriersTablerateHandlingType => Calculate Handling Fee.
// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
var PathCarriersTablerateHandlingType = model.NewStr(`carriers/tablerate/handling_type`)

// PathCarriersTablerateHandlingFee => Handling Fee.
var PathCarriersTablerateHandlingFee = model.NewStr(`carriers/tablerate/handling_fee`)

// PathCarriersTablerateActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersTablerateActive = model.NewBool(`carriers/tablerate/active`)

// PathCarriersTablerateConditionName => Condition.
// SourceModel: Otnegam\OfflineShipping\Model\Config\Source\Tablerate
var PathCarriersTablerateConditionName = model.NewStr(`carriers/tablerate/condition_name`)

// PathCarriersTablerateIncludeVirtualPrice => Include Virtual Products in Price Calculation.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersTablerateIncludeVirtualPrice = model.NewBool(`carriers/tablerate/include_virtual_price`)

// PathCarriersTablerateExport => Export.
var PathCarriersTablerateExport = model.NewStr(`carriers/tablerate/export`)

// PathCarriersTablerateImport => Import.
// BackendModel: Otnegam\OfflineShipping\Model\Config\Backend\Tablerate
var PathCarriersTablerateImport = model.NewStr(`carriers/tablerate/import`)

// PathCarriersTablerateName => Method Name.
var PathCarriersTablerateName = model.NewStr(`carriers/tablerate/name`)

// PathCarriersTablerateSortOrder => Sort Order.
var PathCarriersTablerateSortOrder = model.NewStr(`carriers/tablerate/sort_order`)

// PathCarriersTablerateTitle => Title.
var PathCarriersTablerateTitle = model.NewStr(`carriers/tablerate/title`)

// PathCarriersTablerateSallowspecific => Ship to Applicable Countries.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
var PathCarriersTablerateSallowspecific = model.NewStr(`carriers/tablerate/sallowspecific`)

// PathCarriersTablerateSpecificcountry => Ship to Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathCarriersTablerateSpecificcountry = model.NewStringCSV(`carriers/tablerate/specificcountry`)

// PathCarriersTablerateShowmethod => Show Method if Not Applicable.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersTablerateShowmethod = model.NewBool(`carriers/tablerate/showmethod`)

// PathCarriersTablerateSpecificerrmsg => Displayed Error Message.
var PathCarriersTablerateSpecificerrmsg = model.NewStr(`carriers/tablerate/specificerrmsg`)

// PathCarriersFreeshippingActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersFreeshippingActive = model.NewBool(`carriers/freeshipping/active`)

// PathCarriersFreeshippingFreeShippingSubtotal => Minimum Order Amount.
var PathCarriersFreeshippingFreeShippingSubtotal = model.NewStr(`carriers/freeshipping/free_shipping_subtotal`)

// PathCarriersFreeshippingName => Method Name.
var PathCarriersFreeshippingName = model.NewStr(`carriers/freeshipping/name`)

// PathCarriersFreeshippingSortOrder => Sort Order.
var PathCarriersFreeshippingSortOrder = model.NewStr(`carriers/freeshipping/sort_order`)

// PathCarriersFreeshippingTitle => Title.
var PathCarriersFreeshippingTitle = model.NewStr(`carriers/freeshipping/title`)

// PathCarriersFreeshippingSallowspecific => Ship to Applicable Countries.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
var PathCarriersFreeshippingSallowspecific = model.NewStr(`carriers/freeshipping/sallowspecific`)

// PathCarriersFreeshippingSpecificcountry => Ship to Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathCarriersFreeshippingSpecificcountry = model.NewStringCSV(`carriers/freeshipping/specificcountry`)

// PathCarriersFreeshippingShowmethod => Show Method if Not Applicable.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersFreeshippingShowmethod = model.NewBool(`carriers/freeshipping/showmethod`)

// PathCarriersFreeshippingSpecificerrmsg => Displayed Error Message.
var PathCarriersFreeshippingSpecificerrmsg = model.NewStr(`carriers/freeshipping/specificerrmsg`)
