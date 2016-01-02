// +build ignore

package sitemap

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
	// SitemapCategoryChangefreq => Frequency.
	// Path: sitemap/category/changefreq
	// SourceModel: Otnegam\Sitemap\Model\Config\Source\Frequency
	SitemapCategoryChangefreq model.Str

	// SitemapCategoryPriority => Priority.
	// Valid values range from 0.0 to 1.0.
	// Path: sitemap/category/priority
	// BackendModel: Otnegam\Sitemap\Model\Config\Backend\Priority
	SitemapCategoryPriority model.Str

	// SitemapProductChangefreq => Frequency.
	// Path: sitemap/product/changefreq
	// SourceModel: Otnegam\Sitemap\Model\Config\Source\Frequency
	SitemapProductChangefreq model.Str

	// SitemapProductPriority => Priority.
	// Valid values range from 0.0 to 1.0.
	// Path: sitemap/product/priority
	// BackendModel: Otnegam\Sitemap\Model\Config\Backend\Priority
	SitemapProductPriority model.Str

	// SitemapProductImageInclude => Add Images into Sitemap.
	// Path: sitemap/product/image_include
	// SourceModel: Otnegam\Sitemap\Model\Source\Product\Image\IncludeImage
	SitemapProductImageInclude model.Str

	// SitemapPageChangefreq => Frequency.
	// Path: sitemap/page/changefreq
	// SourceModel: Otnegam\Sitemap\Model\Config\Source\Frequency
	SitemapPageChangefreq model.Str

	// SitemapPagePriority => Priority.
	// Valid values range from 0.0 to 1.0.
	// Path: sitemap/page/priority
	// BackendModel: Otnegam\Sitemap\Model\Config\Backend\Priority
	SitemapPagePriority model.Str

	// SitemapGenerateEnabled => Enabled.
	// Path: sitemap/generate/enabled
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SitemapGenerateEnabled model.Bool

	// SitemapGenerateErrorEmail => Error Email Recipient.
	// Path: sitemap/generate/error_email
	SitemapGenerateErrorEmail model.Str

	// SitemapGenerateErrorEmailIdentity => Error Email Sender.
	// Path: sitemap/generate/error_email_identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	SitemapGenerateErrorEmailIdentity model.Str

	// SitemapGenerateErrorEmailTemplate => Error Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sitemap/generate/error_email_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SitemapGenerateErrorEmailTemplate model.Str

	// SitemapGenerateFrequency => Frequency.
	// Path: sitemap/generate/frequency
	// BackendModel: Otnegam\Cron\Model\Config\Backend\Sitemap
	// SourceModel: Otnegam\Cron\Model\Config\Source\Frequency
	SitemapGenerateFrequency model.Str

	// SitemapGenerateTime => Start Time.
	// Path: sitemap/generate/time
	SitemapGenerateTime model.Str

	// SitemapLimitMaxLines => Maximum No of URLs Per File.
	// Path: sitemap/limit/max_lines
	SitemapLimitMaxLines model.Str

	// SitemapLimitMaxFileSize => Maximum File Size.
	// File size in bytes.
	// Path: sitemap/limit/max_file_size
	SitemapLimitMaxFileSize model.Str

	// SitemapSearchEnginesSubmissionRobots => Enable Submission to Robots.txt.
	// Path: sitemap/search_engines/submission_robots
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SitemapSearchEnginesSubmissionRobots model.Bool
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.SitemapCategoryChangefreq = model.NewStr(`sitemap/category/changefreq`, model.WithPkgCfg(pkgCfg))
	pp.SitemapCategoryPriority = model.NewStr(`sitemap/category/priority`, model.WithPkgCfg(pkgCfg))
	pp.SitemapProductChangefreq = model.NewStr(`sitemap/product/changefreq`, model.WithPkgCfg(pkgCfg))
	pp.SitemapProductPriority = model.NewStr(`sitemap/product/priority`, model.WithPkgCfg(pkgCfg))
	pp.SitemapProductImageInclude = model.NewStr(`sitemap/product/image_include`, model.WithPkgCfg(pkgCfg))
	pp.SitemapPageChangefreq = model.NewStr(`sitemap/page/changefreq`, model.WithPkgCfg(pkgCfg))
	pp.SitemapPagePriority = model.NewStr(`sitemap/page/priority`, model.WithPkgCfg(pkgCfg))
	pp.SitemapGenerateEnabled = model.NewBool(`sitemap/generate/enabled`, model.WithPkgCfg(pkgCfg))
	pp.SitemapGenerateErrorEmail = model.NewStr(`sitemap/generate/error_email`, model.WithPkgCfg(pkgCfg))
	pp.SitemapGenerateErrorEmailIdentity = model.NewStr(`sitemap/generate/error_email_identity`, model.WithPkgCfg(pkgCfg))
	pp.SitemapGenerateErrorEmailTemplate = model.NewStr(`sitemap/generate/error_email_template`, model.WithPkgCfg(pkgCfg))
	pp.SitemapGenerateFrequency = model.NewStr(`sitemap/generate/frequency`, model.WithPkgCfg(pkgCfg))
	pp.SitemapGenerateTime = model.NewStr(`sitemap/generate/time`, model.WithPkgCfg(pkgCfg))
	pp.SitemapLimitMaxLines = model.NewStr(`sitemap/limit/max_lines`, model.WithPkgCfg(pkgCfg))
	pp.SitemapLimitMaxFileSize = model.NewStr(`sitemap/limit/max_file_size`, model.WithPkgCfg(pkgCfg))
	pp.SitemapSearchEnginesSubmissionRobots = model.NewBool(`sitemap/search_engines/submission_robots`, model.WithPkgCfg(pkgCfg))

	return pp
}
