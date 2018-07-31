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

package proto_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/corestoreio/pkg/config"
	cfgval "github.com/corestoreio/pkg/config/validation"
	"github.com/corestoreio/pkg/config/validation/json"
	"github.com/corestoreio/pkg/config/validation/proto"
	"github.com/corestoreio/pkg/util/validation"
	google_protobuf "github.com/gogo/protobuf/types"
	"google.golang.org/grpc"
)

type observerRegistererFake struct {
	*testing.T    // can't be used because of the parallel context of a server.
	wantEvent     uint8
	wantRoute     string
	wantValidator interface{}
	err           error
}

func (orf *observerRegistererFake) RegisterObserver(event uint8, route string, o config.Observer) error {
	if orf.err != nil {
		return orf.err
	}
	if orf.wantEvent != event {
		panic(fmt.Sprintf("Event should be equal: have %d want %d", orf.wantEvent, event))
	}
	if orf.wantRoute != route {
		panic(fmt.Sprintf("Route Want:%q Have:%q", orf.wantRoute, route))
	}

	// Pointers are different in the final objects hence they get printed and
	// their structure compared, not the address.
	if want, have := fmt.Sprintf("%#v", orf.wantValidator), fmt.Sprintf("%#v", o); want != have {
		panic(fmt.Sprintf("Observer internal types should match.\nWant:%s\nHave:%s\n", want, have))
	}
	return nil
}

func (orf *observerRegistererFake) DeregisterObserver(event uint8, route string) error {
	if orf.err != nil {
		return orf.err
	}
	if orf.wantEvent != event {
		panic(fmt.Sprintf("Event should be equal: have %d want %d", orf.wantEvent, event))
	}
	if orf.wantRoute != route {
		panic(fmt.Sprintf("Route Want:%q Have:%q", orf.wantRoute, route))
	}

	return nil
}

const grpcPort = ":50053"

var (
	_ validation.Validator = (*proto.Validator)(nil)
	_ validation.Validator = (*proto.Validators)(nil)
	// _ proto.ConfigValidationServiceServer = (*proto.registrar)(nil)
)

func grpcServer(t *testing.T, or config.ObserverRegisterer, stop <-chan struct{}) {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		t.Fatal(err)
	}
	s := grpc.NewServer()
	proto.RegisterConfigValidationServiceServer(s, proto.NewConfigValidationServiceServer(or))

	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
	go func() {
		<-stop
		s.Stop()
	}()
}

func TestNewConfigValidationServiceServer_Register_Ok(t *testing.T) {
	stop := make(chan struct{})

	grpcServer(t, &observerRegistererFake{
		wantRoute:     "shipment/dhl/free",
		wantEvent:     config.EventOnBeforeSet,
		wantValidator: cfgval.MustNewStrings(cfgval.Strings{Validators: []string{"bool"}}),
	}, stop)

	const address = "localhost" + grpcPort
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer func() {
		stop <- struct{}{}
		assert.NoError(t, conn.Close())
	}()

	c := proto.NewConfigValidationServiceClient(conn)

	result, err := c.Register(context.Background(), &proto.Validators{
		Collection: []*proto.Validator{
			{json.Validator{
				Event:     "before_set",
				Route:     "shipment/dhl/free",
				Type:      "strings",
				Condition: []byte(`{"validators":["bool"]}`),
			}},
		},
	})

	assert.NoError(t, err)
	assert.Exactly(t, &google_protobuf.Empty{}, result)
}

func TestNewConfigValidationServiceServer_Register_Invalid(t *testing.T) {
	stop := make(chan struct{})

	grpcServer(t, &observerRegistererFake{}, stop)

	const address = "localhost" + grpcPort
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer func() {
		stop <- struct{}{}
		assert.NoError(t, conn.Close())
	}()

	c := proto.NewConfigValidationServiceClient(conn)

	result, err := c.Register(context.Background(), &proto.Validators{
		Collection: []*proto.Validator{
			{json.Validator{
				Event:     "before_det",
				Route:     "shipment/dhl/free",
				Type:      "strings",
				Condition: []byte(`{"validators":["bool"]}`),
			}},
		},
	})

	assert.Nil(t, result)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = [config/validation/proto]")
	// TODO further investigate on how to access the original error message
	// st, ok := status.FromError(err)
	// assert.True(t, ok)
	// repr.Println(err)
	// assert.Exactly(t, []interface{}{errors.New("asd")}, st.Details())
}

func TestNewConfigValidationServiceServer_Deegister_Ok(t *testing.T) {
	stop := make(chan struct{})

	grpcServer(t, &observerRegistererFake{
		wantRoute: "shipment/dhl/free",
		wantEvent: config.EventOnBeforeSet,
	}, stop)

	const address = "localhost" + grpcPort
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer func() {
		stop <- struct{}{}
		assert.NoError(t, conn.Close())
	}()

	c := proto.NewConfigValidationServiceClient(conn)

	result, err := c.Deregister(context.Background(), &proto.Validators{
		Collection: []*proto.Validator{
			{json.Validator{
				Event: "before_set",
				Route: "shipment/dhl/free",
			}},
		},
	})

	assert.NoError(t, err)
	assert.Exactly(t, &google_protobuf.Empty{}, result)
}
