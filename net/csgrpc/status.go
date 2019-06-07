// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package csgrpc

import (
	"github.com/gogo/googleapis/google/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewStatusBadRequestError creates a new GRPC status code error with optional
// advanced BadRequest_FieldViolation hints. The argument slice FieldDescription
// must be a balanced slice, where index = Field named and index+1 =
// Description. If msg should parse arguments, use fmt.Sprintf.
func NewStatusBadRequestError(c codes.Code, msg string, FieldDescription ...string) error {
	st := status.New(c, msg)
	if lfd := len(FieldDescription); lfd > 0 && lfd%2 == 0 {
		fvs := make([]*rpc.BadRequest_FieldViolation, lfd/2)
		for i, j := 0, 0; i < lfd; i = i + 2 {
			fvs[j] = &rpc.BadRequest_FieldViolation{
				Field:       FieldDescription[i],
				Description: FieldDescription[i+1],
			}
			j++
		}
		if detSt, err := st.WithDetails(&rpc.BadRequest{FieldViolations: fvs}); err == nil {
			return detSt.Err()
		}
	}

	return st.Err()
}
