// +build ignore

package sitemap

import (
	"github.com/corestoreio/cspkg/config/cfgmodel"
	"github.com/corestoreio/cspkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// SitemapCategoryChangefreq => Frequency.
	// Path: sitemap/category/changefreq
	// SourceModel: Magento\Sitemap\Model\Config\Source\Frequency
	SitemapCategoryChangefreq cfgmodel.Str

	// SitemapCategoryPriority => Priority.
	// Valid values range from 0.0 to 1.0.
	// Path: sitemap/category/priority
	// BackendModel: Magento\Sitemap\Model\Config\Backend\Priority
	SitemapCategoryPriority cfgmodel.Str

	// SitemapProductChangefreq => Frequency.
	// Path: sitemap/product/changefreq
	// SourceModel: Magento\Sitemap\Model\Config\Source\Frequency
	SitemapProductChangefreq cfgmodel.Str

	// SitemapProductPriority => Priority.
	// Valid values range from 0.0 to 1.0.
	// Path: sitemap/product/priority
	// BackendModel: Magento\Sitemap\Model\Config\Backend\Priority
	SitemapProductPriority cfgmodel.Str

	// SitemapProductImageInclude => Add Images into Sitemap.
	// Path: sitemap/product/image_include
	// SourceModel: Magento\Sitemap\Model\Source\Product\Image\IncludeImage
	SitemapProductImageInclude cfgmodel.Str

	// SitemapPageChangefreq => Frequency.
	// Path: sitemap/page/changefreq
	// SourceModel: Magento\Sitemap\Model\Config\Source\Frequency
	SitemapPageChangefreq cfgmodel.Str

	// SitemapPagePriority => Priority.
	// Valid values range from 0.0 to 1.0.
	// Path: sitemap/page/priority
	// BackendModel: Magento\Sitemap\Model\Config\Backend\Priority
	SitemapPagePriority cfgmodel.Str

	// SitemapGenerateEnabled => Enabled.
	// Path: sitemap/generate/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SitemapGenerateEnabled cfgmodel.Bool

	// SitemapGenerateErrorEmail => Error Email Recipient.
	// Path: sitemap/generate/error_email
	SitemapGenerateErrorEmail cfgmodel.Str

	// SitemapGenerateErrorEmailIdentity => Error Email Sender.
	// Path: sitemap/generate/error_email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SitemapGenerateErrorEmailIdentity cfgmodel.Str

	// SitemapGenerateErrorEmailTemplate => Error Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sitemap/generate/error_email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SitemapGenerateErrorEmailTemplate cfgmodel.Str

	// SitemapGenerateFrequency => Frequency.
	// Path: sitemap/generate/frequency
	// BackendModel: Magento\Cron\Model\Config\Backend\Sitemap
	// SourceModel: Magento\Cron\Model\Config\Source\Frequency
	SitemapGenerateFrequency cfgmodel.Str

	// SitemapGenerateTime => Start Time.
	// Path: sitemap/generate/time
	SitemapGenerateTime cfgmodel.Str

	// SitemapLimitMaxLines => Maximum No of URLs Per File.
	// Path: sitemap/limit/max_lines
	SitemapLimitMaxLines cfgmodel.Str

	// SitemapLimitMaxFileSize => Maximum File Size.
	// File size in bytes.
	// Path: sitemap/limit/max_file_size
	SitemapLimitMaxFileSize cfgmodel.Str

	// SitemapSearchEnginesSubmissionRobots => Enable Submission to Robots.txt.
	// Path: sitemap/search_engines/submission_robots
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SitemapSearchEnginesSubmissionRobots cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.SitemapCategoryChangefreq = cfgmodel.NewStr(`sitemap/category/changefreq`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapCategoryPriority = cfgmodel.NewStr(`sitemap/category/priority`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapProductChangefreq = cfgmodel.NewStr(`sitemap/product/changefreq`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapProductPriority = cfgmodel.NewStr(`sitemap/product/priority`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapProductImageInclude = cfgmodel.NewStr(`sitemap/product/image_include`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapPageChangefreq = cfgmodel.NewStr(`sitemap/page/changefreq`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapPagePriority = cfgmodel.NewStr(`sitemap/page/priority`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapGenerateEnabled = cfgmodel.NewBool(`sitemap/generate/enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapGenerateErrorEmail = cfgmodel.NewStr(`sitemap/generate/error_email`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapGenerateErrorEmailIdentity = cfgmodel.NewStr(`sitemap/generate/error_email_identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapGenerateErrorEmailTemplate = cfgmodel.NewStr(`sitemap/generate/error_email_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapGenerateFrequency = cfgmodel.NewStr(`sitemap/generate/frequency`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapGenerateTime = cfgmodel.NewStr(`sitemap/generate/time`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapLimitMaxLines = cfgmodel.NewStr(`sitemap/limit/max_lines`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapLimitMaxFileSize = cfgmodel.NewStr(`sitemap/limit/max_file_size`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapSearchEnginesSubmissionRobots = cfgmodel.NewBool(`sitemap/search_engines/submission_robots`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
