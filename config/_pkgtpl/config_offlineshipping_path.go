// +build ignore

package offlineshipping

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCarriersFlatrateActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersFlatrateActive = model.NewBool(`carriers/flatrate/active`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFlatrateName => Method Name.
var PathCarriersFlatrateName = model.NewStr(`carriers/flatrate/name`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFlatratePrice => Price.
var PathCarriersFlatratePrice = model.NewStr(`carriers/flatrate/price`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFlatrateHandlingType => Calculate Handling Fee.
// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
var PathCarriersFlatrateHandlingType = model.NewStr(`carriers/flatrate/handling_type`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFlatrateHandlingFee => Handling Fee.
var PathCarriersFlatrateHandlingFee = model.NewStr(`carriers/flatrate/handling_fee`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFlatrateSortOrder => Sort Order.
var PathCarriersFlatrateSortOrder = model.NewStr(`carriers/flatrate/sort_order`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFlatrateTitle => Title.
var PathCarriersFlatrateTitle = model.NewStr(`carriers/flatrate/title`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFlatrateType => Type.
// SourceModel: Otnegam\OfflineShipping\Model\Config\Source\Flatrate
var PathCarriersFlatrateType = model.NewStr(`carriers/flatrate/type`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFlatrateSallowspecific => Ship to Applicable Countries.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
var PathCarriersFlatrateSallowspecific = model.NewStr(`carriers/flatrate/sallowspecific`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFlatrateSpecificcountry => Ship to Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathCarriersFlatrateSpecificcountry = model.NewStringCSV(`carriers/flatrate/specificcountry`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFlatrateShowmethod => Show Method if Not Applicable.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersFlatrateShowmethod = model.NewBool(`carriers/flatrate/showmethod`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFlatrateSpecificerrmsg => Displayed Error Message.
var PathCarriersFlatrateSpecificerrmsg = model.NewStr(`carriers/flatrate/specificerrmsg`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersTablerateHandlingType => Calculate Handling Fee.
// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
var PathCarriersTablerateHandlingType = model.NewStr(`carriers/tablerate/handling_type`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersTablerateHandlingFee => Handling Fee.
var PathCarriersTablerateHandlingFee = model.NewStr(`carriers/tablerate/handling_fee`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersTablerateActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersTablerateActive = model.NewBool(`carriers/tablerate/active`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersTablerateConditionName => Condition.
// SourceModel: Otnegam\OfflineShipping\Model\Config\Source\Tablerate
var PathCarriersTablerateConditionName = model.NewStr(`carriers/tablerate/condition_name`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersTablerateIncludeVirtualPrice => Include Virtual Products in Price Calculation.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersTablerateIncludeVirtualPrice = model.NewBool(`carriers/tablerate/include_virtual_price`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersTablerateExport => Export.
var PathCarriersTablerateExport = model.NewStr(`carriers/tablerate/export`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersTablerateImport => Import.
// BackendModel: Otnegam\OfflineShipping\Model\Config\Backend\Tablerate
var PathCarriersTablerateImport = model.NewStr(`carriers/tablerate/import`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersTablerateName => Method Name.
var PathCarriersTablerateName = model.NewStr(`carriers/tablerate/name`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersTablerateSortOrder => Sort Order.
var PathCarriersTablerateSortOrder = model.NewStr(`carriers/tablerate/sort_order`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersTablerateTitle => Title.
var PathCarriersTablerateTitle = model.NewStr(`carriers/tablerate/title`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersTablerateSallowspecific => Ship to Applicable Countries.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
var PathCarriersTablerateSallowspecific = model.NewStr(`carriers/tablerate/sallowspecific`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersTablerateSpecificcountry => Ship to Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathCarriersTablerateSpecificcountry = model.NewStringCSV(`carriers/tablerate/specificcountry`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersTablerateShowmethod => Show Method if Not Applicable.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersTablerateShowmethod = model.NewBool(`carriers/tablerate/showmethod`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersTablerateSpecificerrmsg => Displayed Error Message.
var PathCarriersTablerateSpecificerrmsg = model.NewStr(`carriers/tablerate/specificerrmsg`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFreeshippingActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersFreeshippingActive = model.NewBool(`carriers/freeshipping/active`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFreeshippingFreeShippingSubtotal => Minimum Order Amount.
var PathCarriersFreeshippingFreeShippingSubtotal = model.NewStr(`carriers/freeshipping/free_shipping_subtotal`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFreeshippingName => Method Name.
var PathCarriersFreeshippingName = model.NewStr(`carriers/freeshipping/name`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFreeshippingSortOrder => Sort Order.
var PathCarriersFreeshippingSortOrder = model.NewStr(`carriers/freeshipping/sort_order`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFreeshippingTitle => Title.
var PathCarriersFreeshippingTitle = model.NewStr(`carriers/freeshipping/title`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFreeshippingSallowspecific => Ship to Applicable Countries.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
var PathCarriersFreeshippingSallowspecific = model.NewStr(`carriers/freeshipping/sallowspecific`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFreeshippingSpecificcountry => Ship to Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathCarriersFreeshippingSpecificcountry = model.NewStringCSV(`carriers/freeshipping/specificcountry`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFreeshippingShowmethod => Show Method if Not Applicable.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersFreeshippingShowmethod = model.NewBool(`carriers/freeshipping/showmethod`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersFreeshippingSpecificerrmsg => Displayed Error Message.
var PathCarriersFreeshippingSpecificerrmsg = model.NewStr(`carriers/freeshipping/specificerrmsg`, model.WithPkgCfg(PackageConfiguration))
