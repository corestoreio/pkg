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
	"bytes"
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/utils/mail"
	"github.com/go-gomail/gomail"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

var configMock = config.NewMockReader(
	config.MockInt(func(path string) int {
		println("int", path)
		return 25 // Port 25
	}),
	config.MockString(func(path string) string {
		println("string", path)
		return "localhost"
	}),
	config.MockBool(func(path string) bool {
		//		println("bool", path)
		switch path {
		case "stores/3001/system/smtp/disable":
			return true
		case "stores/4010/system/smtp/disable":
			return false
		default:
			return false
		}
	}),
)

type mockDial struct {
	t        *testing.T
	dial     func()
	dialErr  error
	sendErr  error
	closeErr error
}

func (md mockDial) Dial() (gomail.SendCloser, error) {

	if md.dial != nil {
		md.dial()
	}

	return mockSendCloser{
		t:        md.t,
		sendErr:  md.sendErr,
		closeErr: md.closeErr,
	}, md.dialErr
}

var _ mail.Dialer = (*mockDial)(nil)

type mockSendCloser struct {
	t        *testing.T
	sendErr  error
	closeErr error
}

func (msc mockSendCloser) Send(from string, to []string, msg io.WriterTo) error {
	assert.NotEmpty(msc.t, from)
	assert.NotEmpty(msc.t, to)

	var buf bytes.Buffer
	msg.WriteTo(&buf)
	assert.NotEmpty(msc.t, buf.String())

	//msc.t.Log(buf.String())

	return msc.sendErr
}
func (msc mockSendCloser) Close() error {

	return msc.closeErr
}

var _ gomail.SendCloser = (*mockSendCloser)(nil)

//type mockSender gomail.SendFunc
//
//func (s mockSender) Send(from string, to []string, msg io.WriterTo) error {
//	return s(from, to, msg)
//}
//
//type mockSendCloser struct {
//	mockSender
//	close func() error
//}
//
//func (s *mockSendCloser) Close() error {
//	return s.close()
//}

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
