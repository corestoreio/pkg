// +build ignore

package googleadwords

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
	// GoogleAdwordsActive => Enable.
	// Path: google/adwords/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	GoogleAdwordsActive cfgmodel.Bool

	// GoogleAdwordsConversionId => Conversion ID.
	// Path: google/adwords/conversion_id
	// BackendModel: Magento\GoogleAdwords\Model\Config\Backend\ConversionId
	GoogleAdwordsConversionId cfgmodel.Str

	// GoogleAdwordsConversionLanguage => Conversion Language.
	// Path: google/adwords/conversion_language
	// SourceModel: Magento\GoogleAdwords\Model\Config\Source\Language
	GoogleAdwordsConversionLanguage cfgmodel.Str

	// GoogleAdwordsConversionFormat => Conversion Format.
	// Path: google/adwords/conversion_format
	GoogleAdwordsConversionFormat cfgmodel.Str

	// GoogleAdwordsConversionColor => Conversion Color.
	// Path: google/adwords/conversion_color
	// BackendModel: Magento\GoogleAdwords\Model\Config\Backend\Color
	GoogleAdwordsConversionColor cfgmodel.Str

	// GoogleAdwordsConversionLabel => Conversion Label.
	// Path: google/adwords/conversion_label
	GoogleAdwordsConversionLabel cfgmodel.Str

	// GoogleAdwordsConversionValueType => Conversion Value Type.
	// Path: google/adwords/conversion_value_type
	// SourceModel: Magento\GoogleAdwords\Model\Config\Source\ValueType
	GoogleAdwordsConversionValueType cfgmodel.Str

	// GoogleAdwordsConversionValue => Conversion Value.
	// Path: google/adwords/conversion_value
	GoogleAdwordsConversionValue cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.GoogleAdwordsActive = cfgmodel.NewBool(`google/adwords/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.GoogleAdwordsConversionId = cfgmodel.NewStr(`google/adwords/conversion_id`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.GoogleAdwordsConversionLanguage = cfgmodel.NewStr(`google/adwords/conversion_language`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.GoogleAdwordsConversionFormat = cfgmodel.NewStr(`google/adwords/conversion_format`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.GoogleAdwordsConversionColor = cfgmodel.NewStr(`google/adwords/conversion_color`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.GoogleAdwordsConversionLabel = cfgmodel.NewStr(`google/adwords/conversion_label`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.GoogleAdwordsConversionValueType = cfgmodel.NewStr(`google/adwords/conversion_value_type`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.GoogleAdwordsConversionValue = cfgmodel.NewStr(`google/adwords/conversion_value`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
