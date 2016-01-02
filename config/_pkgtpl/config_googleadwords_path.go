// +build ignore

package googleadwords

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathGoogleAdwordsActive => Enable.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathGoogleAdwordsActive = model.NewBool(`google/adwords/active`, model.WithPkgCfg(PackageConfiguration))

// PathGoogleAdwordsConversionId => Conversion ID.
// BackendModel: Otnegam\GoogleAdwords\Model\Config\Backend\ConversionId
var PathGoogleAdwordsConversionId = model.NewStr(`google/adwords/conversion_id`, model.WithPkgCfg(PackageConfiguration))

// PathGoogleAdwordsConversionLanguage => Conversion Language.
// SourceModel: Otnegam\GoogleAdwords\Model\Config\Source\Language
var PathGoogleAdwordsConversionLanguage = model.NewStr(`google/adwords/conversion_language`, model.WithPkgCfg(PackageConfiguration))

// PathGoogleAdwordsConversionFormat => Conversion Format.
var PathGoogleAdwordsConversionFormat = model.NewStr(`google/adwords/conversion_format`, model.WithPkgCfg(PackageConfiguration))

// PathGoogleAdwordsConversionColor => Conversion Color.
// BackendModel: Otnegam\GoogleAdwords\Model\Config\Backend\Color
var PathGoogleAdwordsConversionColor = model.NewStr(`google/adwords/conversion_color`, model.WithPkgCfg(PackageConfiguration))

// PathGoogleAdwordsConversionLabel => Conversion Label.
var PathGoogleAdwordsConversionLabel = model.NewStr(`google/adwords/conversion_label`, model.WithPkgCfg(PackageConfiguration))

// PathGoogleAdwordsConversionValueType => Conversion Value Type.
// SourceModel: Otnegam\GoogleAdwords\Model\Config\Source\ValueType
var PathGoogleAdwordsConversionValueType = model.NewStr(`google/adwords/conversion_value_type`, model.WithPkgCfg(PackageConfiguration))

// PathGoogleAdwordsConversionValue => Conversion Value.
var PathGoogleAdwordsConversionValue = model.NewStr(`google/adwords/conversion_value`, model.WithPkgCfg(PackageConfiguration))
