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

// +build csall proto

package observer_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/alecthomas/repr"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/observer"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/validation"
	google_protobuf "github.com/gogo/protobuf/types"
	"google.golang.org/grpc"
)

type observerRegistererPanicFake struct {
	testName      string
	wantEvent     uint8
	wantRoute     string
	wantValidator interface{}
	err           error
}

func (orf *observerRegistererPanicFake) RegisterObserver(event uint8, route string, o config.Observer) error {
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
		panic(fmt.Sprintf("%s: Observer internal types should match.\nWant:%s\nHave:%s\n", orf.testName, want, have))
	}
	return nil
}

func (orf *observerRegistererPanicFake) DeregisterObserver(event uint8, route string) error {
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

var (
	_ validation.Validator = (*observer.Configuration)(nil)
	_ validation.Validator = (*observer.Configurations)(nil)
	// _ observer.ProtoServiceServer = (*observer.registrar)(nil)
)

func grpcServer(t *testing.T, port string, or config.ObserverRegisterer, stop <-chan struct{}) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		t.Fatal(err)
	}
	s := grpc.NewServer()
	observer.RegisterProtoServiceServer(s, observer.NewProtoServiceServer(or))

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

func TestNewProtoServiceServer_Register_Ok(t *testing.T) {
	stop := make(chan struct{})

	grpcServer(t, ":50053", &observerRegistererPanicFake{
		testName:      t.Name(),
		wantRoute:     "shipment/dhl/free",
		wantEvent:     config.EventOnBeforeSet,
		wantValidator: observer.MustNewValidator(observer.ValidatorArg{Funcs: []string{"bool"}}),
	}, stop)

	const address = "localhost:50053"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer func() {
		stop <- struct{}{}
		assert.NoError(t, conn.Close())
	}()

	c := observer.NewProtoServiceClient(conn)

	result, err := c.Register(context.Background(), &observer.Configurations{
		Collection: []*observer.Configuration{
			{
				Event:     "before_set",
				Route:     "shipment/dhl/free",
				Type:      "validator",
				Condition: []byte(`{"funcs":["bool"]}`),
			},
		},
	})
	if err != nil {
		t.Log(repr.String(err))
	}
	assert.NoError(t, err)
	assert.Exactly(t, &google_protobuf.Empty{}, result)
}

func TestNewProtoServiceServer_Register_Invalid(t *testing.T) {
	stop := make(chan struct{})

	grpcServer(t, ":50054", &observerRegistererPanicFake{testName: t.Name()}, stop)

	const address = "localhost:50054"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer func() {
		stop <- struct{}{}
		assert.NoError(t, conn.Close())
	}()

	c := observer.NewProtoServiceClient(conn)

	result, err := c.Register(context.Background(), &observer.Configurations{
		Collection: []*observer.Configuration{
			{
				Event:     "before_det",
				Route:     "shipment/dhl/free",
				Type:      "validation",
				Condition: []byte(`{"funcs":["bool"]}`),
			},
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

func TestNewProtoServiceServer_Deegister_Ok(t *testing.T) {
	stop := make(chan struct{})

	grpcServer(t, ":50055", &observerRegistererPanicFake{
		testName:  t.Name(),
		wantRoute: "shipment/dhl/free",
		wantEvent: config.EventOnBeforeSet,
	}, stop)

	const address = "localhost:50055"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer func() {
		stop <- struct{}{}
		assert.NoError(t, conn.Close())
	}()

	c := observer.NewProtoServiceClient(conn)

	result, err := c.Deregister(context.Background(), &observer.Configurations{
		Collection: []*observer.Configuration{
			{
				Event: "before_set",
				Route: "shipment/dhl/free",
			},
		},
	})

	assert.NoError(t, err)
	assert.Exactly(t, &google_protobuf.Empty{}, result)
}
