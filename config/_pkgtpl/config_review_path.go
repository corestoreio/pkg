// +build ignore

package review

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCatalogReviewAllowGuest => Allow Guests to Write Reviews.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogReviewAllowGuest = model.NewBool(`catalog/review/allow_guest`)
