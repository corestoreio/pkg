// +build ignore

package productvideo

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
	// CatalogProductVideoYoutubeApiKey => YouTube API Key.
	// Path: catalog/product_video/youtube_api_key
	CatalogProductVideoYoutubeApiKey model.Str

	// CatalogProductVideoPlayIfBase => Autostart base video.
	// Path: catalog/product_video/play_if_base
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogProductVideoPlayIfBase model.Bool

	// CatalogProductVideoShowRelated => Show related video.
	// Path: catalog/product_video/show_related
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogProductVideoShowRelated model.Bool

	// CatalogProductVideoVideoAutoRestart => Auto restart video.
	// Path: catalog/product_video/video_auto_restart
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogProductVideoVideoAutoRestart model.Bool
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogProductVideoYoutubeApiKey = model.NewStr(`catalog/product_video/youtube_api_key`, model.WithPkgCfg(pkgCfg))
	pp.CatalogProductVideoPlayIfBase = model.NewBool(`catalog/product_video/play_if_base`, model.WithPkgCfg(pkgCfg))
	pp.CatalogProductVideoShowRelated = model.NewBool(`catalog/product_video/show_related`, model.WithPkgCfg(pkgCfg))
	pp.CatalogProductVideoVideoAutoRestart = model.NewBool(`catalog/product_video/video_auto_restart`, model.WithPkgCfg(pkgCfg))

	return pp
}
