// +build ignore

package productvideo

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCatalogProductVideoYoutubeApiKey => YouTube API Key.
var PathCatalogProductVideoYoutubeApiKey = model.NewStr(`catalog/product_video/youtube_api_key`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogProductVideoPlayIfBase => Autostart base video.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogProductVideoPlayIfBase = model.NewBool(`catalog/product_video/play_if_base`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogProductVideoShowRelated => Show related video.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogProductVideoShowRelated = model.NewBool(`catalog/product_video/show_related`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogProductVideoVideoAutoRestart => Auto restart video.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogProductVideoVideoAutoRestart = model.NewBool(`catalog/product_video/video_auto_restart`, model.WithPkgCfg(PackageConfiguration))
