// +build ignore

package sitemap

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathSitemapCategoryChangefreq => Frequency.
// SourceModel: Otnegam\Sitemap\Model\Config\Source\Frequency
var PathSitemapCategoryChangefreq = model.NewStr(`sitemap/category/changefreq`)

// PathSitemapCategoryPriority => Priority.
// Valid values range from 0.0 to 1.0.
// BackendModel: Otnegam\Sitemap\Model\Config\Backend\Priority
var PathSitemapCategoryPriority = model.NewStr(`sitemap/category/priority`)

// PathSitemapProductChangefreq => Frequency.
// SourceModel: Otnegam\Sitemap\Model\Config\Source\Frequency
var PathSitemapProductChangefreq = model.NewStr(`sitemap/product/changefreq`)

// PathSitemapProductPriority => Priority.
// Valid values range from 0.0 to 1.0.
// BackendModel: Otnegam\Sitemap\Model\Config\Backend\Priority
var PathSitemapProductPriority = model.NewStr(`sitemap/product/priority`)

// PathSitemapProductImageInclude => Add Images into Sitemap.
// SourceModel: Otnegam\Sitemap\Model\Source\Product\Image\IncludeImage
var PathSitemapProductImageInclude = model.NewStr(`sitemap/product/image_include`)

// PathSitemapPageChangefreq => Frequency.
// SourceModel: Otnegam\Sitemap\Model\Config\Source\Frequency
var PathSitemapPageChangefreq = model.NewStr(`sitemap/page/changefreq`)

// PathSitemapPagePriority => Priority.
// Valid values range from 0.0 to 1.0.
// BackendModel: Otnegam\Sitemap\Model\Config\Backend\Priority
var PathSitemapPagePriority = model.NewStr(`sitemap/page/priority`)

// PathSitemapGenerateEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSitemapGenerateEnabled = model.NewBool(`sitemap/generate/enabled`)

// PathSitemapGenerateErrorEmail => Error Email Recipient.
var PathSitemapGenerateErrorEmail = model.NewStr(`sitemap/generate/error_email`)

// PathSitemapGenerateErrorEmailIdentity => Error Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSitemapGenerateErrorEmailIdentity = model.NewStr(`sitemap/generate/error_email_identity`)

// PathSitemapGenerateErrorEmailTemplate => Error Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSitemapGenerateErrorEmailTemplate = model.NewStr(`sitemap/generate/error_email_template`)

// PathSitemapGenerateFrequency => Frequency.
// BackendModel: Otnegam\Cron\Model\Config\Backend\Sitemap
// SourceModel: Otnegam\Cron\Model\Config\Source\Frequency
var PathSitemapGenerateFrequency = model.NewStr(`sitemap/generate/frequency`)

// PathSitemapGenerateTime => Start Time.
var PathSitemapGenerateTime = model.NewStr(`sitemap/generate/time`)

// PathSitemapLimitMaxLines => Maximum No of URLs Per File.
var PathSitemapLimitMaxLines = model.NewStr(`sitemap/limit/max_lines`)

// PathSitemapLimitMaxFileSize => Maximum File Size.
// File size in bytes.
var PathSitemapLimitMaxFileSize = model.NewStr(`sitemap/limit/max_file_size`)

// PathSitemapSearchEnginesSubmissionRobots => Enable Submission to Robots.txt.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSitemapSearchEnginesSubmissionRobots = model.NewBool(`sitemap/search_engines/submission_robots`)
