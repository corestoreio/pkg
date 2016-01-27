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

package config_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

var _ config.MessageReceiver = (*testSubscriber)(nil)

type testSubscriber struct {
	t *testing.T
	f func(p path.Path) error
}

func (ts *testSubscriber) MessageConfig(p path.Path) error {
	ts.t.Logf("Message: %s ScopeGroup %s ScopeID %d", p.String(), p.Scope.String(), p.ID)
	return ts.f(p)
}

func TestPubSubBubbling(t *testing.T) {
	defer debugLogBuf.Reset()
	testPath := path.MustNewByParts("aa/bb/cc")

	s := config.NewService()

	_, err := s.Subscribe(path.Route{}, nil)
	assert.EqualError(t, err, path.ErrIncorrectPath.Error())

	subID, err := s.Subscribe(testPath.Route, &testSubscriber{
		t: t,
		f: func(p path.Path) error {
			assert.Equal(t, testPath, p.String())
			if p.Scope == scope.DefaultID {
				assert.Equal(t, int64(0), p.ID)
			} else {
				assert.Equal(t, int64(123), p.ID)
			}
			return nil
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")

	assert.NoError(t, s.Write(config.Value(1), config.Path(testPath), config.Scope(scope.WebsiteID, 123)))
	assert.NoError(t, s.Close())

	t.Log("Before", "testPath", testPath.Route)
	testPath2 := testPath.Clone()
	assert.NoError(t, testPath2.Append(path.NewRoute("Doh")))

	t.Log("After", "testPath", testPath.Route, "testPath2", testPath2.Route)

	// send on closed channel
	assert.NoError(t, s.Write(config.Value(1), config.Path(testPath2), config.Scope(scope.WebsiteID, 3)))
	assert.EqualError(t, s.Close(), config.ErrPublisherClosed.Error())
}

func TestPubSubPanicSimple(t *testing.T) {
	defer debugLogBuf.Reset()
	testPath := "xx/yy/zz"

	s := config.NewService()
	subID, err := s.Subscribe(testPath, &testSubscriber{
		f: func(path string, sg scope.Scope, id int64) error {
			panic("Don't panic!")
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")
	assert.NoError(t, s.Write(config.Value(321), config.Path(testPath), config.ScopeStore(123)))
	assert.NoError(t, s.Close())
	assert.Contains(t, debugLogBuf.String(), `config.pubSub.publish.recover.r recover: "Don't panic!"`)
}

func TestPubSubPanicError(t *testing.T) {
	defer debugLogBuf.Reset()
	testPath := "aa/bb/cc"

	var pErr = errors.New("OMG! Panic!")
	s := config.NewService()
	subID, err := s.Subscribe(testPath, &testSubscriber{
		f: func(path string, sg scope.Scope, id int64) error {
			panic(pErr)
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")
	assert.NoError(t, s.Write(config.Value(321), config.Path(testPath), config.ScopeStore(123)))

	assert.NoError(t, s.Close())
	assert.Contains(t, debugLogBuf.String(), `config.pubSub.publish.recover.err err: OMG! Panic!`)
}

func TestPubSubPanicMultiple(t *testing.T) {
	defer debugLogBuf.Reset()
	s := config.NewService()

	subID, err := s.Subscribe("xx", &testSubscriber{
		f: func(path string, sg scope.Scope, id int64) error {
			assert.Equal(t, "xx/yy/zz", path)
			panic("One: Don't panic!")
		},
	})
	assert.NoError(t, err)
	assert.True(t, subID > 0)

	subID, err = s.Subscribe("xx/yy", &testSubscriber{
		f: func(path string, sg scope.Scope, id int64) error {
			assert.Equal(t, "xx/yy/zz", path)
			panic("Two: Don't panic!")
		},
	})
	assert.NoError(t, err)
	assert.True(t, subID > 0)

	subID, err = s.Subscribe("xx/yy/zz", &testSubscriber{
		f: func(path string, sg scope.Scope, id int64) error {
			assert.Equal(t, "xx/yy/zz", path)
			panic("Three: Don't panic!")
		},
	})
	assert.NoError(t, err)
	assert.True(t, subID > 0)

	assert.NoError(t, s.Write(config.Value(789), config.Path("xx/yy/zz"), config.ScopeStore(987)))
	assert.NoError(t, s.Close())

	assert.Contains(t, debugLogBuf.String(), `config.pubSub.publish.recover.r recover: "One: Don't panic!`)
	assert.Contains(t, debugLogBuf.String(), `config.pubSub.publish.recover.r recover: "Two: Don't panic!"`)
	assert.Contains(t, debugLogBuf.String(), `config.pubSub.publish.recover.r recover: "Three: Don't panic!"`)
}

func TestPubSubUnsubscribe(t *testing.T) {
	defer debugLogBuf.Reset()

	var pErr = errors.New("WTF? Panic!")
	s := config.NewService()
	subID, err := s.Subscribe("xx/yy/zz", &testSubscriber{
		f: func(path string, sg scope.Scope, id int64) error {
			panic(pErr)
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")
	assert.NoError(t, s.Unsubscribe(subID))
	assert.NoError(t, s.Write(config.Value(321), config.Path("xx/yy/zz"), config.ScopeStore(123)))
	assert.NoError(t, s.Close())
	assert.Contains(t, debugLogBuf.String(), `config.Service.Write path: "stores/123/xx/yy/zz" val: 321`)

}

type levelCalls struct {
	sync.Mutex
	level2Calls int
	level3Calls int
}

func TestPubSubEvict(t *testing.T) {
	defer debugLogBuf.Reset()

	levelCall := new(levelCalls)

	var pErr = errors.New("WTF Eviction? Panic!")
	s := config.NewService()
	subID, err := s.Subscribe("xx/yy", &testSubscriber{
		f: func(path string, sg scope.Scope, id int64) error {
			assert.Contains(t, path, "xx/yy")
			// this function gets called 3 times
			levelCall.Lock()
			levelCall.level2Calls++
			levelCall.Unlock()
			return nil
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, subID)

	subID, err = s.Subscribe("xx/yy/zz", &testSubscriber{
		f: func(path string, sg scope.Scope, id int64) error {
			levelCall.Lock()
			levelCall.level3Calls++
			levelCall.Unlock()
			// this function gets called 1 times and then gets removed
			panic(pErr)
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, subID)

	assert.NoError(t, s.Write(config.Value(321), config.Path("xx/yy/zz"), config.ScopeStore(123)))
	assert.NoError(t, s.Write(config.Value(321), config.Path("xx/yy/aa"), config.ScopeStore(123)))
	assert.NoError(t, s.Write(config.Value(321), config.Path("xx/yy/zz"), config.ScopeStore(123)))

	assert.NoError(t, s.Close())

	assert.Contains(t, debugLogBuf.String(), "config.pubSub.publish.recover.err err: WTF Eviction? Panic!")

	levelCall.Lock()
	assert.Equal(t, 3, levelCall.level2Calls)
	assert.Equal(t, 1, levelCall.level3Calls)
	levelCall.Unlock()
	assert.EqualError(t, s.Close(), config.ErrPublisherClosed.Error())
}
