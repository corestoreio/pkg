// +build ignore

package translation

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathDevJsTranslateStrategy => Translation Strategy.
// Please put your store into maintenance mode and redeploy static files after
// changing strategy
// SourceModel: Otnegam\Translation\Model\Js\Config\Source\Strategy
var PathDevJsTranslateStrategy = model.NewStr(`dev/js/translate_strategy`)
