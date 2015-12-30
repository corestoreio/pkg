// +build ignore

package webapi

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathWebapiSoapCharset => Default Response Charset.
// If empty, UTF-8 will be used.
var PathWebapiSoapCharset = model.NewStr(`webapi/soap/charset`)
