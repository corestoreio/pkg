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
	"io"
	std "log"
	"testing"

	"errors"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/utils/log"
	"github.com/corestoreio/csfw/utils/mail"
	"github.com/stretchr/testify/assert"
)

var errLogBuf bytes.Buffer

func init() {
	log.Set(log.NewStdLogger(
		log.StdErrorOption(&errLogBuf, "testErr", std.LstdFlags),
	))
}

func TestDaemonOffline(t *testing.T) {
	offSend := mail.OfflineSend
	defer func() {
		mail.OfflineSend = offSend
	}()

	mail.OfflineSend = func(from string, to []string, msg io.WriterTo) error {
		var buf bytes.Buffer
		_, err := msg.WriteTo(&buf)
		assert.NoError(t, err)
		assert.Equal(t, "gopher@world", from)
		assert.Equal(t, []string{"apple@cupertino"}, to)
		assert.Contains(t, buf.String(), "phoning home")
		assert.Contains(t, buf.String(), "Subject: Phoning home")
		return nil
	}

	dm, err := mail.NewDaemon(
		mail.SetConfig(configMock),
		mail.SetScope(config.ScopeID(3001)),
	)

	assert.NoError(t, err)
	assert.NotNil(t, dm)
	assert.True(t, dm.IsOffline())

	go func() { assert.NoError(t, dm.Worker()) }()
	assert.NoError(t, dm.SendPlain("gopher@world", "apple@cupertino", "Phoning home", "Hey Apple stop phoning home or you become apple puree"))
	assert.NoError(t, dm.Stop())

	assert.EqualError(t, dm.Worker(), mail.ErrMailChannelClosed.Error())
	assert.EqualError(t, dm.Stop(), mail.ErrMailChannelClosed.Error())
	assert.EqualError(t, dm.Send(nil), mail.ErrMailChannelClosed.Error())
	assert.EqualError(t, dm.SendPlain("", "", "", ""), mail.ErrMailChannelClosed.Error())
	assert.EqualError(t, dm.SendHtml("", "", "", ""), mail.ErrMailChannelClosed.Error())
}

func TestDaemonOfflineLogger(t *testing.T) {
	offLog := mail.OfflineLogger
	defer func() {
		mail.OfflineLogger = offLog
	}()

	var logBufI bytes.Buffer
	var logBufE bytes.Buffer
	mail.OfflineLogger = log.NewStdLogger(
		log.StdLevelOption(log.StdLevelInfo),
		log.StdInfoOption(&logBufI, "test", std.LstdFlags),
		log.StdErrorOption(&logBufE, "test", std.LstdFlags),
	)

	dm, err := mail.NewDaemon(
		mail.SetConfig(configMock),
		mail.SetScope(config.ScopeID(3001)),
	)

	assert.NoError(t, err)
	assert.NotNil(t, dm)
	assert.True(t, dm.IsOffline())

	go func() { assert.NoError(t, dm.Worker()) }()
	assert.NoError(t, dm.SendPlain("gopher@earth", "apple@mothership", "Phoning home", "Hey Apple stop phoning home or you become apple puree"))
	assert.NoError(t, dm.Stop())
	assert.True(t, mail.OfflineLogger.IsInfo())

	time.Sleep(time.Millisecond) // waiting for channel to drain

	assert.Contains(t, logBufI.String(), `Send from: "gopher@earth" to: []string{"apple@mothership"} msg: "Mime-Version: 1.0`)
	assert.Empty(t, logBufE.String())

}

func TestDaemonDaemonOptionErrors(t *testing.T) {
	dm, err := mail.NewDaemon(
		mail.SetConfig(nil),
		mail.SetDialer(nil),
		mail.SetSendFunc(nil),
		mail.SetSMTPTimeout(0),
		mail.SetTLSConfig(nil),
		mail.SetScope(nil),
		mail.SetMessageChannel(nil),
	)
	assert.EqualError(t, err, "config.Reader cannot be nil\ngomail.Dialer cannot be nil\ngomail.SendFunc cannot be nil\nTime.Duration cannot be 0\n*tls.Config cannot be nil\nconfig.ScopeIDer cannot be nil\n*gomail.Message channel cannot be nil\n")
	assert.Nil(t, dm)
}

func TestDaemonWorkerDialSend(t *testing.T) {

	dm, err := mail.NewDaemon(
		mail.SetConfig(configMock),
		mail.SetScope(config.ScopeID(4010)),
		mail.SetDialer(
			mockDial{t: t},
		),
	)

	assert.NoError(t, err)
	assert.NotNil(t, dm)
	assert.False(t, dm.IsOffline())

	go func() { assert.NoError(t, dm.Worker()) }()
	assert.NoError(t, dm.SendPlain("rust@lang", "apple@cupertino", "Spagetti", "Pastafari meets Rustafari"))
	assert.NoError(t, dm.Stop())

}

func TestDaemonWorkerDialCloseError(t *testing.T) {
	defer errLogBuf.Reset()
	dm, err := mail.NewDaemon(
		mail.SetConfig(configMock),
		mail.SetSMTPTimeout(time.Millisecond*10),
		mail.SetScope(config.ScopeID(4010)),
		mail.SetDialer(
			mockDial{
				t:        t,
				closeErr: errors.New("Test Close Error"),
			},
		),
	)

	assert.NoError(t, err)
	assert.NotNil(t, dm)
	assert.False(t, dm.IsOffline())

	go func() {
		assert.EqualError(t, dm.Worker(), "Test Close Error", "See goroutine")
	}()
	assert.NoError(t, dm.SendPlain("rust@lang", "apple@cupertino", "Spagetti", "Pastafari meets Rustafari"))
	time.Sleep(time.Millisecond * 100)
	assert.NoError(t, dm.Stop())
	assert.Contains(t, errLogBuf.String(), "mail.daemon.Start.Close err: Test Close Error")

}
