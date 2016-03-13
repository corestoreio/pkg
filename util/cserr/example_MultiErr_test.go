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

package cserr_test

import (
	"errors"
	"fmt"
	"github.com/corestoreio/csfw/util/cserr"
)

type Service struct {
	me *cserr.MultiErr
	a  string
	b  string
}

type Option func(*Service)

func WithA(str string) Option {
	return func(s *Service) {
		if str == "" {
			s.me = s.me.AppendErrors(errors.New("Input string A is empty"))
			return
		}
		s.a = str
	}
}

func WithB(str string) Option {
	return func(s *Service) {
		if str == "" {
			s.me = s.me.AppendErrors(errors.New("Input string B is empty"))
			return
		}
		s.b = str
	}
}

func NewService(opts ...Option) (*Service, error) {
	s := new(Service)
	for _, opt := range opts {
		opt(s)
	}
	if s.me.HasErrors() {
		return nil, s.me
	}
	return s, nil
}

func (s *Service) String() string {
	return s.a
}

func ExampleMultiErr() {

	s1, err := NewService(WithA("hello gophers"))
	fmt.Println("S1:", s1, "; Error is nil:", err == nil)

	s2, err := NewService(WithA(""))
	if err != nil {
		fmt.Println("S2 nil:", s2 == nil, "; S2 Error:", err)
	}

	s3, err := NewService(WithA(""), WithB(""))
	if err != nil {
		fmt.Println("S3 nil:", s3 == nil, "; S3 Error:", err)
	}

	// Output:
	// S1: hello gophers ; Error is nil: true
	// S2 nil: true ; S2 Error: Input string A is empty
	// S3 nil: true ; S3 Error: Input string A is empty
	// Input string B is empty
}
