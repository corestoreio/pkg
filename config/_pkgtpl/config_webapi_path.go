// +build ignore

package webapi

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with ConfigStructure.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// WebapiSoapCharset => Default Response Charset.
	// If empty, UTF-8 will be used.
	// Path: webapi/soap/charset
	WebapiSoapCharset model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.WebapiSoapCharset = model.NewStr(`webapi/soap/charset`, model.WithConfigStructure(cfgStruct))

	return pp
}
