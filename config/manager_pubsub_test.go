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

package config_test

import (
	"bytes"
	std "log"
	"testing"
	"time"

	"errors"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/utils/log"
	"github.com/stretchr/testify/assert"
)

var errLogBuf bytes.Buffer

func init() {
	log.Set(log.NewStdLogger(
		log.SetStdError(&errLogBuf, "testErr: ", std.Lshortfile),
	))
	log.SetLevel(log.StdLevelError)
}

var _ config.Subscriber = (*testSubscriber)(nil)

type testSubscriber struct {
	f func(path string, sg config.ScopeGroup, s config.ScopeIDer)
}

func (ts *testSubscriber) Message(path string, sg config.ScopeGroup, s config.ScopeIDer) {
	//ts.t.Logf("Message: %s Group %d Scope %d", path, sg, s.ScopeID())
	ts.f(path, sg, s)
}

func TestPubSubBubbling(t *testing.T) {

	m := config.NewManager()
	subID, err := m.Subscribe(&testSubscriber{
		f: func(path string, sg config.ScopeGroup, s config.ScopeIDer) {
			assert.Equal(t, "aa/bb/cc", path)
			if sg == config.ScopeDefaultID {
				assert.Equal(t, int64(0), s.ScopeID())
			} else {
				assert.Equal(t, int64(123), s.ScopeID())
			}
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")

	m.Write(config.Value(1), config.Path("aa/bb/cc"), config.Scope(config.ScopeWebsiteID, config.ScopeID(123)))

	assert.NoError(t, m.Close())
	time.Sleep(time.Millisecond * 10) // wait for goroutine to close

	// send on closed channel
	m.Write(config.NoBubble(), config.Value(1), config.Path("a/b/c"), config.Scope(config.ScopeWebsiteID, config.ScopeID(3)))
	assert.EqualError(t, m.Close(), config.ErrPublisherClosed.Error())
}

func TestPubSubPanic(t *testing.T) {
	defer errLogBuf.Reset()
	m := config.NewManager()
	subID, err := m.Subscribe(&testSubscriber{
		f: func(path string, sg config.ScopeGroup, s config.ScopeIDer) {
			panic("Don't panic!")
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")
	m.Write(config.Value(321), config.Path("x/y/z"), config.ScopeStore(config.ScopeID(123)), config.NoBubble())
	assert.NoError(t, m.Close())
	time.Sleep(time.Millisecond * 10) // wait for goroutine to close
	assert.Contains(t, errLogBuf.String(), `config.pubSub.publish.recover.wtf recover: "Don't panic!"`)
}

func TestPubSubPanicError(t *testing.T) {
	defer errLogBuf.Reset()
	var pErr = errors.New("OMG! Panic!")
	m := config.NewManager()
	subID, err := m.Subscribe(&testSubscriber{
		f: func(path string, sg config.ScopeGroup, s config.ScopeIDer) {
			panic(pErr)
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")
	m.Write(config.Value(321), config.Path("x/y/z"), config.ScopeStore(config.ScopeID(123)), config.NoBubble())
	// not closing channel to let the Goroutine around egging aka. herumeiern.
	time.Sleep(time.Millisecond * 10) // wait for goroutine ...
	assert.Contains(t, errLogBuf.String(), `config.pubSub.publish.recover.err err: OMG! Panic!`)
}

func TestPubSubPanicMultiple(t *testing.T) {
	defer errLogBuf.Reset()
	m := config.NewManager()

	subID, err := m.Subscribe(&testSubscriber{
		f: func(path string, sg config.ScopeGroup, s config.ScopeIDer) {
			panic("One: Don't panic!")
		},
	})
	assert.NoError(t, err)
	assert.True(t, subID > 0)

	subID, err = m.Subscribe(&testSubscriber{
		f: func(path string, sg config.ScopeGroup, s config.ScopeIDer) {
			panic("Two: Don't panic!")
		},
	})
	assert.NoError(t, err)
	assert.True(t, subID > 0)

	subID, err = m.Subscribe(&testSubscriber{
		f: func(path string, sg config.ScopeGroup, s config.ScopeIDer) {
			panic("Three: Don't panic!")
		},
	})
	assert.NoError(t, err)
	assert.True(t, subID > 0)

	m.Write(config.Value(789), config.Path("x/y/z"), config.ScopeStore(config.ScopeID(987)), config.NoBubble())
	assert.NoError(t, m.Close())
	time.Sleep(time.Millisecond * 30) // wait for goroutine to close
	assert.Contains(t, errLogBuf.String(), `testErr: stdLib.go:228: config.pubSub.publish.recover.wtf recover: "One: Don't panic!"
testErr: stdLib.go:228: config.pubSub.publish.recover.wtf recover: "Two: Don't panic!"
testErr: stdLib.go:228: config.pubSub.publish.recover.wtf recover: "Three: Don't panic!"`+"\n")
}

func TestPubSubUnsubscribe(t *testing.T) {
	defer errLogBuf.Reset()

	var pErr = errors.New("WTF? Panic!")
	m := config.NewManager()
	subID, err := m.Subscribe(&testSubscriber{
		f: func(path string, sg config.ScopeGroup, s config.ScopeIDer) {
			panic(pErr)
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")
	assert.NoError(t, m.Unsubscribe(subID))
	m.Write(config.Value(321), config.Path("x/y/z"), config.ScopeStore(config.ScopeID(123)), config.NoBubble())
	// not closing channel to let the Goroutine around egging aka. herumeiern.
	time.Sleep(time.Millisecond) // wait for goroutine ...
	assert.Empty(t, errLogBuf.String())
}
