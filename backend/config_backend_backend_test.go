package backend_test

import (
	"github.com/corestoreio/csfw/backend"
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/config/source"
	"testing"
)

// benchmarkGlobalStruct trick the compiler to not optimize anything
var benchmarkGlobalStruct bool

func Benchmark_StructGlobal(b *testing.B) {

	sg := config.NewMockGetter().NewScoped(1, 1, 1)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkGlobalStruct, err = backend.Backend.DevCSSMinifyFiles.Get(sg) // any random struct field
		if err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_StructSpecific(b *testing.B) {

	sg := config.NewMockGetter().NewScoped(1, 1, 1)

	mb := model.NewBool("aa/bb/cc", model.WithConfigStructure(backend.ConfigStructure), model.WithSource(source.YesNo))

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkGlobalStruct, err = mb.Get(sg)
		if err != nil {
			b.Error(err)
		}
	}
}
