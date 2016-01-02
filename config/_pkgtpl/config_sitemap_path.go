// +build ignore

package sitemap

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathSitemapCategoryChangefreq => Frequency.
// SourceModel: Otnegam\Sitemap\Model\Config\Source\Frequency
var PathSitemapCategoryChangefreq = model.NewStr(`sitemap/category/changefreq`, model.WithPkgCfg(PackageConfiguration))

// PathSitemapCategoryPriority => Priority.
// Valid values range from 0.0 to 1.0.
// BackendModel: Otnegam\Sitemap\Model\Config\Backend\Priority
var PathSitemapCategoryPriority = model.NewStr(`sitemap/category/priority`, model.WithPkgCfg(PackageConfiguration))

// PathSitemapProductChangefreq => Frequency.
// SourceModel: Otnegam\Sitemap\Model\Config\Source\Frequency
var PathSitemapProductChangefreq = model.NewStr(`sitemap/product/changefreq`, model.WithPkgCfg(PackageConfiguration))

// PathSitemapProductPriority => Priority.
// Valid values range from 0.0 to 1.0.
// BackendModel: Otnegam\Sitemap\Model\Config\Backend\Priority
var PathSitemapProductPriority = model.NewStr(`sitemap/product/priority`, model.WithPkgCfg(PackageConfiguration))

// PathSitemapProductImageInclude => Add Images into Sitemap.
// SourceModel: Otnegam\Sitemap\Model\Source\Product\Image\IncludeImage
var PathSitemapProductImageInclude = model.NewStr(`sitemap/product/image_include`, model.WithPkgCfg(PackageConfiguration))

// PathSitemapPageChangefreq => Frequency.
// SourceModel: Otnegam\Sitemap\Model\Config\Source\Frequency
var PathSitemapPageChangefreq = model.NewStr(`sitemap/page/changefreq`, model.WithPkgCfg(PackageConfiguration))

// PathSitemapPagePriority => Priority.
// Valid values range from 0.0 to 1.0.
// BackendModel: Otnegam\Sitemap\Model\Config\Backend\Priority
var PathSitemapPagePriority = model.NewStr(`sitemap/page/priority`, model.WithPkgCfg(PackageConfiguration))

// PathSitemapGenerateEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSitemapGenerateEnabled = model.NewBool(`sitemap/generate/enabled`, model.WithPkgCfg(PackageConfiguration))

// PathSitemapGenerateErrorEmail => Error Email Recipient.
var PathSitemapGenerateErrorEmail = model.NewStr(`sitemap/generate/error_email`, model.WithPkgCfg(PackageConfiguration))

// PathSitemapGenerateErrorEmailIdentity => Error Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSitemapGenerateErrorEmailIdentity = model.NewStr(`sitemap/generate/error_email_identity`, model.WithPkgCfg(PackageConfiguration))

// PathSitemapGenerateErrorEmailTemplate => Error Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSitemapGenerateErrorEmailTemplate = model.NewStr(`sitemap/generate/error_email_template`, model.WithPkgCfg(PackageConfiguration))

// PathSitemapGenerateFrequency => Frequency.
// BackendModel: Otnegam\Cron\Model\Config\Backend\Sitemap
// SourceModel: Otnegam\Cron\Model\Config\Source\Frequency
var PathSitemapGenerateFrequency = model.NewStr(`sitemap/generate/frequency`, model.WithPkgCfg(PackageConfiguration))

// PathSitemapGenerateTime => Start Time.
var PathSitemapGenerateTime = model.NewStr(`sitemap/generate/time`, model.WithPkgCfg(PackageConfiguration))

// PathSitemapLimitMaxLines => Maximum No of URLs Per File.
var PathSitemapLimitMaxLines = model.NewStr(`sitemap/limit/max_lines`, model.WithPkgCfg(PackageConfiguration))

// PathSitemapLimitMaxFileSize => Maximum File Size.
// File size in bytes.
var PathSitemapLimitMaxFileSize = model.NewStr(`sitemap/limit/max_file_size`, model.WithPkgCfg(PackageConfiguration))

// PathSitemapSearchEnginesSubmissionRobots => Enable Submission to Robots.txt.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSitemapSearchEnginesSubmissionRobots = model.NewBool(`sitemap/search_engines/submission_robots`, model.WithPkgCfg(PackageConfiguration))
