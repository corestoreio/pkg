// +build ignore

package offlineshipping

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	model.PkgBackend
	// CarriersFlatrateActive => Enabled.
	// Path: carriers/flatrate/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFlatrateActive model.Bool

	// CarriersFlatrateName => Method Name.
	// Path: carriers/flatrate/name
	CarriersFlatrateName model.Str

	// CarriersFlatratePrice => Price.
	// Path: carriers/flatrate/price
	CarriersFlatratePrice model.Str

	// CarriersFlatrateHandlingType => Calculate Handling Fee.
	// Path: carriers/flatrate/handling_type
	// SourceModel: Magento\Shipping\Model\Source\HandlingType
	CarriersFlatrateHandlingType model.Str

	// CarriersFlatrateHandlingFee => Handling Fee.
	// Path: carriers/flatrate/handling_fee
	CarriersFlatrateHandlingFee model.Str

	// CarriersFlatrateSortOrder => Sort Order.
	// Path: carriers/flatrate/sort_order
	CarriersFlatrateSortOrder model.Str

	// CarriersFlatrateTitle => Title.
	// Path: carriers/flatrate/title
	CarriersFlatrateTitle model.Str

	// CarriersFlatrateType => Type.
	// Path: carriers/flatrate/type
	// SourceModel: Magento\OfflineShipping\Model\Config\Source\Flatrate
	CarriersFlatrateType model.Str

	// CarriersFlatrateSallowspecific => Ship to Applicable Countries.
	// Path: carriers/flatrate/sallowspecific
	// SourceModel: Magento\Shipping\Model\Config\Source\Allspecificcountries
	CarriersFlatrateSallowspecific model.Str

	// CarriersFlatrateSpecificcountry => Ship to Specific Countries.
	// Path: carriers/flatrate/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	CarriersFlatrateSpecificcountry model.StringCSV

	// CarriersFlatrateShowmethod => Show Method if Not Applicable.
	// Path: carriers/flatrate/showmethod
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFlatrateShowmethod model.Bool

	// CarriersFlatrateSpecificerrmsg => Displayed Error Message.
	// Path: carriers/flatrate/specificerrmsg
	CarriersFlatrateSpecificerrmsg model.Str

	// CarriersTablerateHandlingType => Calculate Handling Fee.
	// Path: carriers/tablerate/handling_type
	// SourceModel: Magento\Shipping\Model\Source\HandlingType
	CarriersTablerateHandlingType model.Str

	// CarriersTablerateHandlingFee => Handling Fee.
	// Path: carriers/tablerate/handling_fee
	CarriersTablerateHandlingFee model.Str

	// CarriersTablerateActive => Enabled.
	// Path: carriers/tablerate/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersTablerateActive model.Bool

	// CarriersTablerateConditionName => Condition.
	// Path: carriers/tablerate/condition_name
	// SourceModel: Magento\OfflineShipping\Model\Config\Source\Tablerate
	CarriersTablerateConditionName model.Str

	// CarriersTablerateIncludeVirtualPrice => Include Virtual Products in Price Calculation.
	// Path: carriers/tablerate/include_virtual_price
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersTablerateIncludeVirtualPrice model.Bool

	// CarriersTablerateExport => Export.
	// Path: carriers/tablerate/export
	CarriersTablerateExport model.Str

	// CarriersTablerateImport => Import.
	// Path: carriers/tablerate/import
	// BackendModel: Magento\OfflineShipping\Model\Config\Backend\Tablerate
	CarriersTablerateImport model.Str

	// CarriersTablerateName => Method Name.
	// Path: carriers/tablerate/name
	CarriersTablerateName model.Str

	// CarriersTablerateSortOrder => Sort Order.
	// Path: carriers/tablerate/sort_order
	CarriersTablerateSortOrder model.Str

	// CarriersTablerateTitle => Title.
	// Path: carriers/tablerate/title
	CarriersTablerateTitle model.Str

	// CarriersTablerateSallowspecific => Ship to Applicable Countries.
	// Path: carriers/tablerate/sallowspecific
	// SourceModel: Magento\Shipping\Model\Config\Source\Allspecificcountries
	CarriersTablerateSallowspecific model.Str

	// CarriersTablerateSpecificcountry => Ship to Specific Countries.
	// Path: carriers/tablerate/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	CarriersTablerateSpecificcountry model.StringCSV

	// CarriersTablerateShowmethod => Show Method if Not Applicable.
	// Path: carriers/tablerate/showmethod
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersTablerateShowmethod model.Bool

	// CarriersTablerateSpecificerrmsg => Displayed Error Message.
	// Path: carriers/tablerate/specificerrmsg
	CarriersTablerateSpecificerrmsg model.Str

	// CarriersFreeshippingActive => Enabled.
	// Path: carriers/freeshipping/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFreeshippingActive model.Bool

	// CarriersFreeshippingFreeShippingSubtotal => Minimum Order Amount.
	// Path: carriers/freeshipping/free_shipping_subtotal
	CarriersFreeshippingFreeShippingSubtotal model.Str

	// CarriersFreeshippingName => Method Name.
	// Path: carriers/freeshipping/name
	CarriersFreeshippingName model.Str

	// CarriersFreeshippingSortOrder => Sort Order.
	// Path: carriers/freeshipping/sort_order
	CarriersFreeshippingSortOrder model.Str

	// CarriersFreeshippingTitle => Title.
	// Path: carriers/freeshipping/title
	CarriersFreeshippingTitle model.Str

	// CarriersFreeshippingSallowspecific => Ship to Applicable Countries.
	// Path: carriers/freeshipping/sallowspecific
	// SourceModel: Magento\Shipping\Model\Config\Source\Allspecificcountries
	CarriersFreeshippingSallowspecific model.Str

	// CarriersFreeshippingSpecificcountry => Ship to Specific Countries.
	// Path: carriers/freeshipping/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	CarriersFreeshippingSpecificcountry model.StringCSV

	// CarriersFreeshippingShowmethod => Show Method if Not Applicable.
	// Path: carriers/freeshipping/showmethod
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFreeshippingShowmethod model.Bool

	// CarriersFreeshippingSpecificerrmsg => Displayed Error Message.
	// Path: carriers/freeshipping/specificerrmsg
	CarriersFreeshippingSpecificerrmsg model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CarriersFlatrateActive = model.NewBool(`carriers/flatrate/active`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFlatrateName = model.NewStr(`carriers/flatrate/name`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFlatratePrice = model.NewStr(`carriers/flatrate/price`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFlatrateHandlingType = model.NewStr(`carriers/flatrate/handling_type`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFlatrateHandlingFee = model.NewStr(`carriers/flatrate/handling_fee`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFlatrateSortOrder = model.NewStr(`carriers/flatrate/sort_order`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFlatrateTitle = model.NewStr(`carriers/flatrate/title`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFlatrateType = model.NewStr(`carriers/flatrate/type`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFlatrateSallowspecific = model.NewStr(`carriers/flatrate/sallowspecific`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFlatrateSpecificcountry = model.NewStringCSV(`carriers/flatrate/specificcountry`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFlatrateShowmethod = model.NewBool(`carriers/flatrate/showmethod`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFlatrateSpecificerrmsg = model.NewStr(`carriers/flatrate/specificerrmsg`, model.WithConfigStructure(cfgStruct))
	pp.CarriersTablerateHandlingType = model.NewStr(`carriers/tablerate/handling_type`, model.WithConfigStructure(cfgStruct))
	pp.CarriersTablerateHandlingFee = model.NewStr(`carriers/tablerate/handling_fee`, model.WithConfigStructure(cfgStruct))
	pp.CarriersTablerateActive = model.NewBool(`carriers/tablerate/active`, model.WithConfigStructure(cfgStruct))
	pp.CarriersTablerateConditionName = model.NewStr(`carriers/tablerate/condition_name`, model.WithConfigStructure(cfgStruct))
	pp.CarriersTablerateIncludeVirtualPrice = model.NewBool(`carriers/tablerate/include_virtual_price`, model.WithConfigStructure(cfgStruct))
	pp.CarriersTablerateExport = model.NewStr(`carriers/tablerate/export`, model.WithConfigStructure(cfgStruct))
	pp.CarriersTablerateImport = model.NewStr(`carriers/tablerate/import`, model.WithConfigStructure(cfgStruct))
	pp.CarriersTablerateName = model.NewStr(`carriers/tablerate/name`, model.WithConfigStructure(cfgStruct))
	pp.CarriersTablerateSortOrder = model.NewStr(`carriers/tablerate/sort_order`, model.WithConfigStructure(cfgStruct))
	pp.CarriersTablerateTitle = model.NewStr(`carriers/tablerate/title`, model.WithConfigStructure(cfgStruct))
	pp.CarriersTablerateSallowspecific = model.NewStr(`carriers/tablerate/sallowspecific`, model.WithConfigStructure(cfgStruct))
	pp.CarriersTablerateSpecificcountry = model.NewStringCSV(`carriers/tablerate/specificcountry`, model.WithConfigStructure(cfgStruct))
	pp.CarriersTablerateShowmethod = model.NewBool(`carriers/tablerate/showmethod`, model.WithConfigStructure(cfgStruct))
	pp.CarriersTablerateSpecificerrmsg = model.NewStr(`carriers/tablerate/specificerrmsg`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFreeshippingActive = model.NewBool(`carriers/freeshipping/active`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFreeshippingFreeShippingSubtotal = model.NewStr(`carriers/freeshipping/free_shipping_subtotal`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFreeshippingName = model.NewStr(`carriers/freeshipping/name`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFreeshippingSortOrder = model.NewStr(`carriers/freeshipping/sort_order`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFreeshippingTitle = model.NewStr(`carriers/freeshipping/title`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFreeshippingSallowspecific = model.NewStr(`carriers/freeshipping/sallowspecific`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFreeshippingSpecificcountry = model.NewStringCSV(`carriers/freeshipping/specificcountry`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFreeshippingShowmethod = model.NewBool(`carriers/freeshipping/showmethod`, model.WithConfigStructure(cfgStruct))
	pp.CarriersFreeshippingSpecificerrmsg = model.NewStr(`carriers/freeshipping/specificerrmsg`, model.WithConfigStructure(cfgStruct))

	return pp
}
