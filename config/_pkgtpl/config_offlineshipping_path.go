// +build ignore

package offlineshipping

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with PackageConfiguration.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// CarriersFlatrateActive => Enabled.
	// Path: carriers/flatrate/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersFlatrateActive model.Bool

	// CarriersFlatrateName => Method Name.
	// Path: carriers/flatrate/name
	CarriersFlatrateName model.Str

	// CarriersFlatratePrice => Price.
	// Path: carriers/flatrate/price
	CarriersFlatratePrice model.Str

	// CarriersFlatrateHandlingType => Calculate Handling Fee.
	// Path: carriers/flatrate/handling_type
	// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
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
	// SourceModel: Otnegam\OfflineShipping\Model\Config\Source\Flatrate
	CarriersFlatrateType model.Str

	// CarriersFlatrateSallowspecific => Ship to Applicable Countries.
	// Path: carriers/flatrate/sallowspecific
	// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
	CarriersFlatrateSallowspecific model.Str

	// CarriersFlatrateSpecificcountry => Ship to Specific Countries.
	// Path: carriers/flatrate/specificcountry
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	CarriersFlatrateSpecificcountry model.StringCSV

	// CarriersFlatrateShowmethod => Show Method if Not Applicable.
	// Path: carriers/flatrate/showmethod
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersFlatrateShowmethod model.Bool

	// CarriersFlatrateSpecificerrmsg => Displayed Error Message.
	// Path: carriers/flatrate/specificerrmsg
	CarriersFlatrateSpecificerrmsg model.Str

	// CarriersTablerateHandlingType => Calculate Handling Fee.
	// Path: carriers/tablerate/handling_type
	// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
	CarriersTablerateHandlingType model.Str

	// CarriersTablerateHandlingFee => Handling Fee.
	// Path: carriers/tablerate/handling_fee
	CarriersTablerateHandlingFee model.Str

	// CarriersTablerateActive => Enabled.
	// Path: carriers/tablerate/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersTablerateActive model.Bool

	// CarriersTablerateConditionName => Condition.
	// Path: carriers/tablerate/condition_name
	// SourceModel: Otnegam\OfflineShipping\Model\Config\Source\Tablerate
	CarriersTablerateConditionName model.Str

	// CarriersTablerateIncludeVirtualPrice => Include Virtual Products in Price Calculation.
	// Path: carriers/tablerate/include_virtual_price
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersTablerateIncludeVirtualPrice model.Bool

	// CarriersTablerateExport => Export.
	// Path: carriers/tablerate/export
	CarriersTablerateExport model.Str

	// CarriersTablerateImport => Import.
	// Path: carriers/tablerate/import
	// BackendModel: Otnegam\OfflineShipping\Model\Config\Backend\Tablerate
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
	// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
	CarriersTablerateSallowspecific model.Str

	// CarriersTablerateSpecificcountry => Ship to Specific Countries.
	// Path: carriers/tablerate/specificcountry
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	CarriersTablerateSpecificcountry model.StringCSV

	// CarriersTablerateShowmethod => Show Method if Not Applicable.
	// Path: carriers/tablerate/showmethod
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersTablerateShowmethod model.Bool

	// CarriersTablerateSpecificerrmsg => Displayed Error Message.
	// Path: carriers/tablerate/specificerrmsg
	CarriersTablerateSpecificerrmsg model.Str

	// CarriersFreeshippingActive => Enabled.
	// Path: carriers/freeshipping/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
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
	// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
	CarriersFreeshippingSallowspecific model.Str

	// CarriersFreeshippingSpecificcountry => Ship to Specific Countries.
	// Path: carriers/freeshipping/specificcountry
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	CarriersFreeshippingSpecificcountry model.StringCSV

	// CarriersFreeshippingShowmethod => Show Method if Not Applicable.
	// Path: carriers/freeshipping/showmethod
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersFreeshippingShowmethod model.Bool

	// CarriersFreeshippingSpecificerrmsg => Displayed Error Message.
	// Path: carriers/freeshipping/specificerrmsg
	CarriersFreeshippingSpecificerrmsg model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CarriersFlatrateActive = model.NewBool(`carriers/flatrate/active`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFlatrateName = model.NewStr(`carriers/flatrate/name`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFlatratePrice = model.NewStr(`carriers/flatrate/price`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFlatrateHandlingType = model.NewStr(`carriers/flatrate/handling_type`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFlatrateHandlingFee = model.NewStr(`carriers/flatrate/handling_fee`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFlatrateSortOrder = model.NewStr(`carriers/flatrate/sort_order`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFlatrateTitle = model.NewStr(`carriers/flatrate/title`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFlatrateType = model.NewStr(`carriers/flatrate/type`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFlatrateSallowspecific = model.NewStr(`carriers/flatrate/sallowspecific`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFlatrateSpecificcountry = model.NewStringCSV(`carriers/flatrate/specificcountry`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFlatrateShowmethod = model.NewBool(`carriers/flatrate/showmethod`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFlatrateSpecificerrmsg = model.NewStr(`carriers/flatrate/specificerrmsg`, model.WithPkgCfg(pkgCfg))
	pp.CarriersTablerateHandlingType = model.NewStr(`carriers/tablerate/handling_type`, model.WithPkgCfg(pkgCfg))
	pp.CarriersTablerateHandlingFee = model.NewStr(`carriers/tablerate/handling_fee`, model.WithPkgCfg(pkgCfg))
	pp.CarriersTablerateActive = model.NewBool(`carriers/tablerate/active`, model.WithPkgCfg(pkgCfg))
	pp.CarriersTablerateConditionName = model.NewStr(`carriers/tablerate/condition_name`, model.WithPkgCfg(pkgCfg))
	pp.CarriersTablerateIncludeVirtualPrice = model.NewBool(`carriers/tablerate/include_virtual_price`, model.WithPkgCfg(pkgCfg))
	pp.CarriersTablerateExport = model.NewStr(`carriers/tablerate/export`, model.WithPkgCfg(pkgCfg))
	pp.CarriersTablerateImport = model.NewStr(`carriers/tablerate/import`, model.WithPkgCfg(pkgCfg))
	pp.CarriersTablerateName = model.NewStr(`carriers/tablerate/name`, model.WithPkgCfg(pkgCfg))
	pp.CarriersTablerateSortOrder = model.NewStr(`carriers/tablerate/sort_order`, model.WithPkgCfg(pkgCfg))
	pp.CarriersTablerateTitle = model.NewStr(`carriers/tablerate/title`, model.WithPkgCfg(pkgCfg))
	pp.CarriersTablerateSallowspecific = model.NewStr(`carriers/tablerate/sallowspecific`, model.WithPkgCfg(pkgCfg))
	pp.CarriersTablerateSpecificcountry = model.NewStringCSV(`carriers/tablerate/specificcountry`, model.WithPkgCfg(pkgCfg))
	pp.CarriersTablerateShowmethod = model.NewBool(`carriers/tablerate/showmethod`, model.WithPkgCfg(pkgCfg))
	pp.CarriersTablerateSpecificerrmsg = model.NewStr(`carriers/tablerate/specificerrmsg`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFreeshippingActive = model.NewBool(`carriers/freeshipping/active`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFreeshippingFreeShippingSubtotal = model.NewStr(`carriers/freeshipping/free_shipping_subtotal`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFreeshippingName = model.NewStr(`carriers/freeshipping/name`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFreeshippingSortOrder = model.NewStr(`carriers/freeshipping/sort_order`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFreeshippingTitle = model.NewStr(`carriers/freeshipping/title`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFreeshippingSallowspecific = model.NewStr(`carriers/freeshipping/sallowspecific`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFreeshippingSpecificcountry = model.NewStringCSV(`carriers/freeshipping/specificcountry`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFreeshippingShowmethod = model.NewBool(`carriers/freeshipping/showmethod`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFreeshippingSpecificerrmsg = model.NewStr(`carriers/freeshipping/specificerrmsg`, model.WithPkgCfg(pkgCfg))

	return pp
}
