// +build csall proto

package urlvalues

// ProtoToValues converts a proto message to a Values type. It appends to
// argument vals, which can be nil.
func ProtoToValues(vals Values, pkv *ProtoKeyValues) Values {
	if vals == nil {
		vals = make(Values, len(pkv.Data))
	}
	for _, kv := range pkv.Data {
		vals[kv.Key] = kv.Value
	}
	return vals
}

// ValuesToProto converts a Values map to a proto message.
func ValuesToProto(vals Values) *ProtoKeyValues {
	var pkv ProtoKeyValues
	pkv.Data = make([]*ProtoKeyValue, 0, len(vals))
	for k, v := range vals {
		pkv.Data = append(pkv.Data, &ProtoKeyValue{
			Key:   k,
			Value: v,
		})
	}
	return &pkv
}
