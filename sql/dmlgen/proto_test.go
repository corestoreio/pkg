package dmlgen

import (
	"testing"

	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/codegen"
)

func TestBuildDMLProtoMapper(t *testing.T) {
	err := buildDMLProtoMapper(
		"github.com/corestoreio/pkg/sql/dmlgen/dmltestgenerated",
		"github.com/corestoreio/pkg/sql/dmlgen/dmltestgenerated/dmltestgeneratedpb",
	)
	if fe, ok := err.(*codegen.FormatError); ok {
		t.Log(fe.Error())
		t.Log(fe.Code)
	} else {
		assert.NoError(t, err)
	}
}
