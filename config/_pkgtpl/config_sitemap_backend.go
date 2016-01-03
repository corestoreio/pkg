// +build ignore

package sitemap

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	model.PkgBackend
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

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.SitemapCategoryChangefreq = model.NewStr(`sitemap/category/changefreq`, model.WithConfigStructure(cfgStruct))
	pp.SitemapCategoryPriority = model.NewStr(`sitemap/category/priority`, model.WithConfigStructure(cfgStruct))
	pp.SitemapProductChangefreq = model.NewStr(`sitemap/product/changefreq`, model.WithConfigStructure(cfgStruct))
	pp.SitemapProductPriority = model.NewStr(`sitemap/product/priority`, model.WithConfigStructure(cfgStruct))
	pp.SitemapProductImageInclude = model.NewStr(`sitemap/product/image_include`, model.WithConfigStructure(cfgStruct))
	pp.SitemapPageChangefreq = model.NewStr(`sitemap/page/changefreq`, model.WithConfigStructure(cfgStruct))
	pp.SitemapPagePriority = model.NewStr(`sitemap/page/priority`, model.WithConfigStructure(cfgStruct))
	pp.SitemapGenerateEnabled = model.NewBool(`sitemap/generate/enabled`, model.WithConfigStructure(cfgStruct))
	pp.SitemapGenerateErrorEmail = model.NewStr(`sitemap/generate/error_email`, model.WithConfigStructure(cfgStruct))
	pp.SitemapGenerateErrorEmailIdentity = model.NewStr(`sitemap/generate/error_email_identity`, model.WithConfigStructure(cfgStruct))
	pp.SitemapGenerateErrorEmailTemplate = model.NewStr(`sitemap/generate/error_email_template`, model.WithConfigStructure(cfgStruct))
	pp.SitemapGenerateFrequency = model.NewStr(`sitemap/generate/frequency`, model.WithConfigStructure(cfgStruct))
	pp.SitemapGenerateTime = model.NewStr(`sitemap/generate/time`, model.WithConfigStructure(cfgStruct))
	pp.SitemapLimitMaxLines = model.NewStr(`sitemap/limit/max_lines`, model.WithConfigStructure(cfgStruct))
	pp.SitemapLimitMaxFileSize = model.NewStr(`sitemap/limit/max_file_size`, model.WithConfigStructure(cfgStruct))
	pp.SitemapSearchEnginesSubmissionRobots = model.NewBool(`sitemap/search_engines/submission_robots`, model.WithConfigStructure(cfgStruct))

	return pp
}
