// UnmarshalJSON implements interface json.Unmarshaler.
func (cc *{{$.Collection}}) UnmarshalJSON(b []byte) (err error) {
	return json.Unmarshal(b, cc.Data)
}

// MarshalJSON implements interface json.Marshaler.
func (cc *{{$.Collection}}) MarshalJSON() ([]byte, error) {
	return json.Marshal(cc.Data)
}

// TODO add MarshalText and UnmarshalText.
