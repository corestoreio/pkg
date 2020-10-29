//go:generate go get -u github.com/jteeuwen/go-bindata/go-bindata
//go:generate go install github.com/jteeuwen/go-bindata/go-bindata
//go:generate go-bindata -ignore \.git\S* -ignore LICENSE -ignore README\.md -ignore blns\.base64\.txt -ignore blns\.txt -ignore package\.json -o internal/resource.go -nocompress -pkg internal ..

// Package naughtystrings is a collection of strings that have a high
// probability of causing issues when used as user input.
package naughtystrings

import (
	"encoding/json"
	"sync"

	"github.com/corestoreio/pkg/util/naughtystrings/internal"
)

var base64encoded, unencoded []string

var loadOnce sync.Once

func init() {
	loadOnce.Do(func() {
		if base64encoded == nil {
			base64encoded = load("../blns.base64.json")
		}
		if unencoded == nil {
			unencoded = load("../blns.json")
		}
	})
}

// Base64Encoded returns the strings encoded in base 64.
func Base64Encoded() []string {
	return base64encoded
}

// Unencoded returns the strings.
func Unencoded() []string {
	return unencoded
}

func load(file string) []string {
	asset, err := internal.Asset(file)
	if err != nil {
		panic(err)
	}

	var naughty []string

	if err := json.Unmarshal(asset, &naughty); err != nil {
		panic(err)
	}

	return naughty
}
