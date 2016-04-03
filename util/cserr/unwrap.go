// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cserr

import "github.com/juju/errors"

// UnwrapMasked returns the underlying original error to be used for
// comparison with Err* variables. Recursive function.
// Only in use with github.com/juju/errors package.
func UnwrapMasked(err error) error {
	if err == nil {
		return nil
	}
	if theErr, ok := err.(*errors.Err); ok {
		if uErr := theErr.Underlying(); uErr != nil {
			return UnwrapMasked(uErr)
		}
	}
	return err
}
