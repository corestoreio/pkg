// +build ignore

package offlineshipping

import (
	"github.com/corestoreio/cspkg/config/cfgmodel"
	"github.com/corestoreio/cspkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// CarriersFlatrateActive => Enabled.
	// Path: carriers/flatrate/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFlatrateActive cfgmodel.Bool

	// CarriersFlatrateName => Method Name.
	// Path: carriers/flatrate/name
	CarriersFlatrateName cfgmodel.Str

	// CarriersFlatratePrice => Price.
	// Path: carriers/flatrate/price
	CarriersFlatratePrice cfgmodel.Str

	// CarriersFlatrateHandlingType => Calculate Handling Fee.
	// Path: carriers/flatrate/handling_type
	// SourceModel: Magento\Shipping\Model\Source\HandlingType
	CarriersFlatrateHandlingType cfgmodel.Str

	// CarriersFlatrateHandlingFee => Handling Fee.
	// Path: carriers/flatrate/handling_fee
	CarriersFlatrateHandlingFee cfgmodel.Str

	// CarriersFlatrateSortOrder => Sort Order.
	// Path: carriers/flatrate/sort_order
	CarriersFlatrateSortOrder cfgmodel.Str

	// CarriersFlatrateTitle => Title.
	// Path: carriers/flatrate/title
	CarriersFlatrateTitle cfgmodel.Str

	// CarriersFlatrateType => Type.
	// Path: carriers/flatrate/type
	// SourceModel: Magento\OfflineShipping\Model\Config\Source\Flatrate
	CarriersFlatrateType cfgmodel.Str

	// CarriersFlatrateSallowspecific => Ship to Applicable Countries.
	// Path: carriers/flatrate/sallowspecific
	// SourceModel: Magento\Shipping\Model\Config\Source\Allspecificcountries
	CarriersFlatrateSallowspecific cfgmodel.Str

	// CarriersFlatrateSpecificcountry => Ship to Specific Countries.
	// Path: carriers/flatrate/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	CarriersFlatrateSpecificcountry cfgmodel.StringCSV

	// CarriersFlatrateShowmethod => Show Method if Not Applicable.
	// Path: carriers/flatrate/showmethod
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFlatrateShowmethod cfgmodel.Bool

	// CarriersFlatrateSpecificerrmsg => Displayed Error Message.
	// Path: carriers/flatrate/specificerrmsg
	CarriersFlatrateSpecificerrmsg cfgmodel.Str

	// CarriersTablerateHandlingType => Calculate Handling Fee.
	// Path: carriers/tablerate/handling_type
	// SourceModel: Magento\Shipping\Model\Source\HandlingType
	CarriersTablerateHandlingType cfgmodel.Str

	// CarriersTablerateHandlingFee => Handling Fee.
	// Path: carriers/tablerate/handling_fee
	CarriersTablerateHandlingFee cfgmodel.Str

	// CarriersTablerateActive => Enabled.
	// Path: carriers/tablerate/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersTablerateActive cfgmodel.Bool

	// CarriersTablerateConditionName => Condition.
	// Path: carriers/tablerate/condition_name
	// SourceModel: Magento\OfflineShipping\Model\Config\Source\Tablerate
	CarriersTablerateConditionName cfgmodel.Str

	// CarriersTablerateIncludeVirtualPrice => Include Virtual Products in Price Calculation.
	// Path: carriers/tablerate/include_virtual_price
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersTablerateIncludeVirtualPrice cfgmodel.Bool

	// CarriersTablerateExport => Export.
	// Path: carriers/tablerate/export
	CarriersTablerateExport cfgmodel.Str

	// CarriersTablerateImport => Import.
	// Path: carriers/tablerate/import
	// BackendModel: Magento\OfflineShipping\Model\Config\Backend\Tablerate
	CarriersTablerateImport cfgmodel.Str

	// CarriersTablerateName => Method Name.
	// Path: carriers/tablerate/name
	CarriersTablerateName cfgmodel.Str

	// CarriersTablerateSortOrder => Sort Order.
	// Path: carriers/tablerate/sort_order
	CarriersTablerateSortOrder cfgmodel.Str

	// CarriersTablerateTitle => Title.
	// Path: carriers/tablerate/title
	CarriersTablerateTitle cfgmodel.Str

	// CarriersTablerateSallowspecific => Ship to Applicable Countries.
	// Path: carriers/tablerate/sallowspecific
	// SourceModel: Magento\Shipping\Model\Config\Source\Allspecificcountries
	CarriersTablerateSallowspecific cfgmodel.Str

	// CarriersTablerateSpecificcountry => Ship to Specific Countries.
	// Path: carriers/tablerate/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	CarriersTablerateSpecificcountry cfgmodel.StringCSV

	// CarriersTablerateShowmethod => Show Method if Not Applicable.
	// Path: carriers/tablerate/showmethod
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersTablerateShowmethod cfgmodel.Bool

	// CarriersTablerateSpecificerrmsg => Displayed Error Message.
	// Path: carriers/tablerate/specificerrmsg
	CarriersTablerateSpecificerrmsg cfgmodel.Str

	// CarriersFreeshippingActive => Enabled.
	// Path: carriers/freeshipping/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFreeshippingActive cfgmodel.Bool

	// CarriersFreeshippingFreeShippingSubtotal => Minimum Order Amount.
	// Path: carriers/freeshipping/free_shipping_subtotal
	CarriersFreeshippingFreeShippingSubtotal cfgmodel.Str

	// CarriersFreeshippingName => Method Name.
	// Path: carriers/freeshipping/name
	CarriersFreeshippingName cfgmodel.Str

	// CarriersFreeshippingSortOrder => Sort Order.
	// Path: carriers/freeshipping/sort_order
	CarriersFreeshippingSortOrder cfgmodel.Str

	// CarriersFreeshippingTitle => Title.
	// Path: carriers/freeshipping/title
	CarriersFreeshippingTitle cfgmodel.Str

	// CarriersFreeshippingSallowspecific => Ship to Applicable Countries.
	// Path: carriers/freeshipping/sallowspecific
	// SourceModel: Magento\Shipping\Model\Config\Source\Allspecificcountries
	CarriersFreeshippingSallowspecific cfgmodel.Str

	// CarriersFreeshippingSpecificcountry => Ship to Specific Countries.
	// Path: carriers/freeshipping/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	CarriersFreeshippingSpecificcountry cfgmodel.StringCSV

	// CarriersFreeshippingShowmethod => Show Method if Not Applicable.
	// Path: carriers/freeshipping/showmethod
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFreeshippingShowmethod cfgmodel.Bool

	// CarriersFreeshippingSpecificerrmsg => Displayed Error Message.
	// Path: carriers/freeshipping/specificerrmsg
	CarriersFreeshippingSpecificerrmsg cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CarriersFlatrateActive = cfgmodel.NewBool(`carriers/flatrate/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFlatrateName = cfgmodel.NewStr(`carriers/flatrate/name`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFlatratePrice = cfgmodel.NewStr(`carriers/flatrate/price`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFlatrateHandlingType = cfgmodel.NewStr(`carriers/flatrate/handling_type`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFlatrateHandlingFee = cfgmodel.NewStr(`carriers/flatrate/handling_fee`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFlatrateSortOrder = cfgmodel.NewStr(`carriers/flatrate/sort_order`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFlatrateTitle = cfgmodel.NewStr(`carriers/flatrate/title`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFlatrateType = cfgmodel.NewStr(`carriers/flatrate/type`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFlatrateSallowspecific = cfgmodel.NewStr(`carriers/flatrate/sallowspecific`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFlatrateSpecificcountry = cfgmodel.NewStringCSV(`carriers/flatrate/specificcountry`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFlatrateShowmethod = cfgmodel.NewBool(`carriers/flatrate/showmethod`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFlatrateSpecificerrmsg = cfgmodel.NewStr(`carriers/flatrate/specificerrmsg`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersTablerateHandlingType = cfgmodel.NewStr(`carriers/tablerate/handling_type`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersTablerateHandlingFee = cfgmodel.NewStr(`carriers/tablerate/handling_fee`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersTablerateActive = cfgmodel.NewBool(`carriers/tablerate/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersTablerateConditionName = cfgmodel.NewStr(`carriers/tablerate/condition_name`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersTablerateIncludeVirtualPrice = cfgmodel.NewBool(`carriers/tablerate/include_virtual_price`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersTablerateExport = cfgmodel.NewStr(`carriers/tablerate/export`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersTablerateImport = cfgmodel.NewStr(`carriers/tablerate/import`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersTablerateName = cfgmodel.NewStr(`carriers/tablerate/name`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersTablerateSortOrder = cfgmodel.NewStr(`carriers/tablerate/sort_order`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersTablerateTitle = cfgmodel.NewStr(`carriers/tablerate/title`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersTablerateSallowspecific = cfgmodel.NewStr(`carriers/tablerate/sallowspecific`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersTablerateSpecificcountry = cfgmodel.NewStringCSV(`carriers/tablerate/specificcountry`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersTablerateShowmethod = cfgmodel.NewBool(`carriers/tablerate/showmethod`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersTablerateSpecificerrmsg = cfgmodel.NewStr(`carriers/tablerate/specificerrmsg`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFreeshippingActive = cfgmodel.NewBool(`carriers/freeshipping/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFreeshippingFreeShippingSubtotal = cfgmodel.NewStr(`carriers/freeshipping/free_shipping_subtotal`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFreeshippingName = cfgmodel.NewStr(`carriers/freeshipping/name`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFreeshippingSortOrder = cfgmodel.NewStr(`carriers/freeshipping/sort_order`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFreeshippingTitle = cfgmodel.NewStr(`carriers/freeshipping/title`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFreeshippingSallowspecific = cfgmodel.NewStr(`carriers/freeshipping/sallowspecific`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFreeshippingSpecificcountry = cfgmodel.NewStringCSV(`carriers/freeshipping/specificcountry`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFreeshippingShowmethod = cfgmodel.NewBool(`carriers/freeshipping/showmethod`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFreeshippingSpecificerrmsg = cfgmodel.NewStr(`carriers/freeshipping/specificerrmsg`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
