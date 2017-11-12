// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (cc *{{$.Collection}}) UnmarshalBinary(data []byte) error {
	return cc.Unmarshal(data) // Implemented via github.com/gogo/protobuf
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (cc *{{$.Collection}}) MarshalBinary() (data []byte, err error) {
	return cc.Marshal()  // Implemented via github.com/gogo/protobuf
}

// GobDecode kept for Go 1 compatibility reasons.
// deprecated in Go 2, use UnmarshalBinary
func (cc *{{$.Collection}}) GobDecode(data []byte) error {
	return cc.Unmarshal(data) // Implemented via github.com/gogo/protobuf
}

// GobEncode kept for Go 1 compatibility reasons.
// deprecated in Go 2, use MarshalBinary
func (cc *{{$.Collection}}) GobEncode() ([]byte, error) {
	return cc.Marshal()  // Implemented via github.com/gogo/protobuf
}
