// +build ignore

package googleoptimizer

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathGoogleAnalyticsExperiments => Enable Content Experiments.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathGoogleAnalyticsExperiments = model.NewBool(`google/analytics/experiments`, model.WithPkgCfg(PackageConfiguration))
