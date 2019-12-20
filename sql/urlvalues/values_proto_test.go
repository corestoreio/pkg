// +build csall proto

package urlvalues

import (
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestProtoToValues(t *testing.T) {
	vals := ProtoToValues(nil, &ProtoKeyValues{
		Data: []*ProtoKeyValue{
			{Key: "a", Value: []string{"b"}},
		},
	})
	assert.Exactly(t, Values{"a": []string{"b"}}, vals)
}

func TestValuesToProto(t *testing.T) {
	pkv := ValuesToProto(Values{"a": []string{"b"}})
	assert.Exactly(t, &ProtoKeyValues{
		Data: []*ProtoKeyValue{
			{Key: "a", Value: []string{"b"}},
		},
	}, pkv)
}
