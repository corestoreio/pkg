// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (cc *{{$.Collection}}) UnmarshalBinary(data []byte) error {
	return cc.Unmarshal(data) // Implemented via github.com/gogo/protobuf
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (cc *{{$.Collection}}) MarshalBinary() (data []byte, err error) {
	return cc.Marshal()  // Implemented via github.com/gogo/protobuf
}
