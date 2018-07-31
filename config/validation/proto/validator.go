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

package proto

import (
	"context"
	"fmt"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/validation/json"
	"github.com/corestoreio/pkg/net/cspb"
	"github.com/gogo/protobuf/types"
	"google.golang.org/grpc/codes"
)

// Validator defines the data retrieved from the outside as JSON to add a new
// validator for a specific route and event.
type Validator struct {
	json.Validator
}

// Validators a list of Validator types.
type Validators struct {
	Collection []*Validator
}

// Validate validates the current slice and used in GRPC middleware.
func (m Validators) Validate() error {
	for _, v := range m.Collection {
		if err := v.Validate(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

type registrar struct {
	or config.ObserverRegisterer
}

// NewConfigValidationServiceServer creates a new GRPC server which can register
// and deregister validation observers to the concrete config.Server type.
func NewConfigValidationServiceServer(or config.ObserverRegisterer) ConfigValidationServiceServer {
	return registrar{
		or: or,
	}
}

func (r registrar) Register(ctx context.Context, vs *Validators) (*types.Empty, error) {
	if err := vs.Validate(); err != nil {
		return nil, cspb.NewStatusBadRequestError(codes.InvalidArgument, "[config/validation/proto]",
			"Validators.Validate",
			err.Error(),
		)
	}

	for idx, v := range vs.Collection {
		event, route, o, err := v.MakeObserver()
		if err != nil {
			return nil, cspb.NewStatusBadRequestError(codes.InvalidArgument, "[config/validation/proto]",
				fmt.Sprintf("validator_%d", idx),
				err.Error(),
			)
		}
		if err := r.or.RegisterObserver(event, route, o); err != nil {
			return nil, cspb.NewStatusBadRequestError(codes.Internal, "[config/validation/proto]",
				fmt.Sprintf("validator_%d", idx),
				err.Error(),
				"event",
				fmt.Sprintf("%d", event),
				"route",
				route,
			)
		}
	}

	return &types.Empty{}, nil
}

func (r registrar) Deregister(ctx context.Context, vs *Validators) (*types.Empty, error) {
	for _, v := range vs.Collection {
		event, route, err := v.MakeEventRoute()
		if err != nil {
			return nil, errors.Wrapf(err, "[config/validation] Data: %#v", v)
		}
		if err := r.or.DeregisterObserver(event, route); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return &types.Empty{}, nil
}
