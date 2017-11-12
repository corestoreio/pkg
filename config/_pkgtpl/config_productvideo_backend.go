// +build ignore

package productvideo

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
	// CatalogProductVideoYoutubeApiKey => YouTube API Key.
	// Path: catalog/product_video/youtube_api_key
	CatalogProductVideoYoutubeApiKey cfgmodel.Str

	// CatalogProductVideoPlayIfBase => Autostart base video.
	// Path: catalog/product_video/play_if_base
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogProductVideoPlayIfBase cfgmodel.Bool

	// CatalogProductVideoShowRelated => Show related video.
	// Path: catalog/product_video/show_related
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogProductVideoShowRelated cfgmodel.Bool

	// CatalogProductVideoVideoAutoRestart => Auto restart video.
	// Path: catalog/product_video/video_auto_restart
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogProductVideoVideoAutoRestart cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogProductVideoYoutubeApiKey = cfgmodel.NewStr(`catalog/product_video/youtube_api_key`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogProductVideoPlayIfBase = cfgmodel.NewBool(`catalog/product_video/play_if_base`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogProductVideoShowRelated = cfgmodel.NewBool(`catalog/product_video/show_related`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogProductVideoVideoAutoRestart = cfgmodel.NewBool(`catalog/product_video/video_auto_restart`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
