// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package mail_test

import (
	"io"
	"testing"

	"github.com/go-gomail/gomail"
)

type mockSender gomail.SendFunc

func (s mockSender) Send(from string, to []string, msg io.WriterTo) error {
	return s(from, to, msg)
}

type mockSendCloser struct {
	mockSender
	close func() error
}

func (s *mockSendCloser) Close() error {
	return s.close()
}

func TestDaemon(t *testing.T) {

}

//type mockTransport struct {
//	rt func(req *http.Request) (resp *http.Response, err error)
//}
//
//func (t *mockTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
//	return t.rt(req)
//}
//
//func main() {
//	fmt.Println("Hello, playground")
//	tr := &mockTransport{
//		rt: func(r *http.Request) (w *http.Response, err error) {
//			return nil, errors.New("no response")
//		},
//	}
//	c := &http.Client{Transport: tr}
//	resp, err := c.Get("http://github.com/")
//
//	fmt.Println(resp, err)
//}
