package csjwt

import "time"

// Claimer for a type to be a Claims object
type Claimer interface {
	// Valid method that determines if the token is invalid for any supported reason.
	// Returns nil on success
	Valid() error
	// Expires declares when a token expires. A duration smaller or equal
	// to zero means that the token has already expired.
	// Useful when adding a token to a blacklist.
	Expires() time.Duration
	// Set sets a value to the claim and may overwrite existing values
	Set(key string, value interface{}) error
	// Get retrieves a value from the claim.
	Get(key string) (value interface{}, err error)
}

// (CS) I personally don't like the Set() and Get() functions but there is no
// other way around it.
