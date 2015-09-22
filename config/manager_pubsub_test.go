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
	"github.com/corestoreio/csfw/config/scope"
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

var _ config.MessageReceiver = (*testSubscriber)(nil)

type testSubscriber struct {
	f func(path string, sg scope.Group, s scope.IDer) error
}

func (ts *testSubscriber) MessageConfig(path string, sg scope.Group, s scope.IDer) error {
	//ts.t.Logf("Message: %s Group %d Scope %d", path, sg, s.scope.ID())
	return ts.f(path, sg, s)
}

func TestPubSubBubbling(t *testing.T) {
	defer errLogBuf.Reset()
	testPath := "a/b/c"

	m := config.NewManager()

	_, err := m.Subscribe("", nil)
	assert.EqualError(t, err, config.ErrPathEmpty.Error())

	subID, err := m.Subscribe(testPath, &testSubscriber{
		f: func(path string, sg scope.Group, s scope.IDer) error {
			assert.Equal(t, testPath, path)
			if sg == scope.DefaultID {
				assert.Equal(t, int64(0), s.ScopeID())
			} else {
				assert.Equal(t, int64(123), s.ScopeID())
			}
			return nil
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")

	m.Write(config.Value(1), config.Path(testPath), config.Scope(scope.WebsiteID, scope.ID(123)))

	assert.NoError(t, m.Close())
	time.Sleep(time.Millisecond * 10) // wait for goroutine to close

	// send on closed channel
	m.Write(config.NoBubble(), config.Value(1), config.Path(testPath+"Doh"), config.Scope(scope.WebsiteID, scope.ID(3)))
	assert.EqualError(t, m.Close(), config.ErrPublisherClosed.Error())
}

func TestPubSubPanic(t *testing.T) {
	defer errLogBuf.Reset()
	testPath := "x/y/z"

	m := config.NewManager()
	subID, err := m.Subscribe(testPath, &testSubscriber{
		f: func(path string, sg scope.Group, s scope.IDer) error {
			panic("Don't panic!")
			return nil
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")
	m.Write(config.Value(321), config.Path(testPath), config.ScopeStore(scope.ID(123)), config.NoBubble())
	assert.NoError(t, m.Close())
	time.Sleep(time.Millisecond * 10) // wait for goroutine to close
	assert.Contains(t, errLogBuf.String(), `config.pubSub.publish.recover.r recover: "Don't panic!"`)
}

func TestPubSubPanicError(t *testing.T) {
	defer errLogBuf.Reset()
	testPath := "™/ö/º"

	var pErr = errors.New("OMG! Panic!")
	m := config.NewManager()
	subID, err := m.Subscribe(testPath, &testSubscriber{
		f: func(path string, sg scope.Group, s scope.IDer) error {
			panic(pErr)
			return nil
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")
	m.Write(config.Value(321), config.Path(testPath), config.ScopeStore(scope.ID(123)), config.NoBubble())
	// not closing channel to let the Goroutine around egging aka. herumeiern.
	time.Sleep(time.Millisecond * 10) // wait for goroutine ...
	assert.Contains(t, errLogBuf.String(), `config.pubSub.publish.recover.err err: OMG! Panic!`)
}

func TestPubSubPanicMultiple(t *testing.T) {
	defer errLogBuf.Reset()
	m := config.NewManager()

	subID, err := m.Subscribe("x", &testSubscriber{
		f: func(path string, sg scope.Group, s scope.IDer) error {
			assert.Equal(t, "x/y/z", path)
			panic("One: Don't panic!")
			return nil
		},
	})
	assert.NoError(t, err)
	assert.True(t, subID > 0)

	subID, err = m.Subscribe("x/y", &testSubscriber{
		f: func(path string, sg scope.Group, s scope.IDer) error {
			assert.Equal(t, "x/y/z", path)
			panic("Two: Don't panic!")
			return nil
		},
	})
	assert.NoError(t, err)
	assert.True(t, subID > 0)

	subID, err = m.Subscribe("x/y/z", &testSubscriber{
		f: func(path string, sg scope.Group, s scope.IDer) error {
			assert.Equal(t, "x/y/z", path)
			panic("Three: Don't panic!")
			return nil
		},
	})
	assert.NoError(t, err)
	assert.True(t, subID > 0)

	m.Write(config.Value(789), config.Path("x/y/z"), config.ScopeStore(scope.ID(987)), config.NoBubble())
	assert.NoError(t, m.Close())
	time.Sleep(time.Millisecond * 30) // wait for goroutine to close
	assert.Contains(t, errLogBuf.String(), `testErr: stdLib.go:228: config.pubSub.publish.recover.r recover: "One: Don't panic!`)
	assert.Contains(t, errLogBuf.String(), `testErr: stdLib.go:228: config.pubSub.publish.recover.r recover: "Two: Don't panic!"`)
	assert.Contains(t, errLogBuf.String(), `testErr: stdLib.go:228: config.pubSub.publish.recover.r recover: "Three: Don't panic!"`)
}

func TestPubSubUnsubscribe(t *testing.T) {
	defer errLogBuf.Reset()

	var pErr = errors.New("WTF? Panic!")
	m := config.NewManager()
	subID, err := m.Subscribe("x/y/z", &testSubscriber{
		f: func(path string, sg scope.Group, s scope.IDer) error {
			panic(pErr)
			return nil
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")
	assert.NoError(t, m.Unsubscribe(subID))
	m.Write(config.Value(321), config.Path("x/y/z"), config.ScopeStore(scope.ID(123)), config.NoBubble())
	time.Sleep(time.Millisecond) // wait for goroutine ...
	assert.Empty(t, errLogBuf.String())
}

func TestPubSubEvict(t *testing.T) {
	defer errLogBuf.Reset()

	var level2Calls int
	var level3Calls int

	var pErr = errors.New("WTF Eviction? Panic!")
	m := config.NewManager()
	subID, err := m.Subscribe("x/y", &testSubscriber{
		f: func(path string, sg scope.Group, s scope.IDer) error {
			assert.Contains(t, path, "x/y")
			// this function gets called 3 times
			level2Calls++
			return nil
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID)

	subID, err = m.Subscribe("x/y/z", &testSubscriber{
		f: func(path string, sg scope.Group, s scope.IDer) error {
			level3Calls++
			// this function gets called 1 times and then gets removed
			panic(pErr)
			return nil
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, subID)

	m.Write(config.Value(321), config.Path("x/y/z"), config.ScopeStore(scope.ID(123)), config.NoBubble())
	m.Write(config.Value(321), config.Path("x/y/a"), config.ScopeStore(scope.ID(123)), config.NoBubble())
	m.Write(config.Value(321), config.Path("x/y/z"), config.ScopeStore(scope.ID(123)), config.NoBubble())

	time.Sleep(time.Millisecond * 20) // wait for goroutine ...

	assert.Contains(t, errLogBuf.String(), "testErr: stdLib.go:228: config.pubSub.publish.recover.err err: WTF Eviction? Panic!")

	assert.Equal(t, 3, level2Calls)
	assert.Equal(t, 1, level3Calls)
}
