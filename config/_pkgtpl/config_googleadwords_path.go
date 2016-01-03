// +build ignore

package googleadwords

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with ConfigStructure.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// GoogleAdwordsActive => Enable.
	// Path: google/adwords/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	GoogleAdwordsActive model.Bool

	// GoogleAdwordsConversionId => Conversion ID.
	// Path: google/adwords/conversion_id
	// BackendModel: Otnegam\GoogleAdwords\Model\Config\Backend\ConversionId
	GoogleAdwordsConversionId model.Str

	// GoogleAdwordsConversionLanguage => Conversion Language.
	// Path: google/adwords/conversion_language
	// SourceModel: Otnegam\GoogleAdwords\Model\Config\Source\Language
	GoogleAdwordsConversionLanguage model.Str

	// GoogleAdwordsConversionFormat => Conversion Format.
	// Path: google/adwords/conversion_format
	GoogleAdwordsConversionFormat model.Str

	// GoogleAdwordsConversionColor => Conversion Color.
	// Path: google/adwords/conversion_color
	// BackendModel: Otnegam\GoogleAdwords\Model\Config\Backend\Color
	GoogleAdwordsConversionColor model.Str

	// GoogleAdwordsConversionLabel => Conversion Label.
	// Path: google/adwords/conversion_label
	GoogleAdwordsConversionLabel model.Str

	// GoogleAdwordsConversionValueType => Conversion Value Type.
	// Path: google/adwords/conversion_value_type
	// SourceModel: Otnegam\GoogleAdwords\Model\Config\Source\ValueType
	GoogleAdwordsConversionValueType model.Str

	// GoogleAdwordsConversionValue => Conversion Value.
	// Path: google/adwords/conversion_value
	GoogleAdwordsConversionValue model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.GoogleAdwordsActive = model.NewBool(`google/adwords/active`, model.WithConfigStructure(cfgStruct))
	pp.GoogleAdwordsConversionId = model.NewStr(`google/adwords/conversion_id`, model.WithConfigStructure(cfgStruct))
	pp.GoogleAdwordsConversionLanguage = model.NewStr(`google/adwords/conversion_language`, model.WithConfigStructure(cfgStruct))
	pp.GoogleAdwordsConversionFormat = model.NewStr(`google/adwords/conversion_format`, model.WithConfigStructure(cfgStruct))
	pp.GoogleAdwordsConversionColor = model.NewStr(`google/adwords/conversion_color`, model.WithConfigStructure(cfgStruct))
	pp.GoogleAdwordsConversionLabel = model.NewStr(`google/adwords/conversion_label`, model.WithConfigStructure(cfgStruct))
	pp.GoogleAdwordsConversionValueType = model.NewStr(`google/adwords/conversion_value_type`, model.WithConfigStructure(cfgStruct))
	pp.GoogleAdwordsConversionValue = model.NewStr(`google/adwords/conversion_value`, model.WithConfigStructure(cfgStruct))

	return pp
}
