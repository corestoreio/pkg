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

package config_test

// TODO use github.com/fortytw2/leaktest

// import (
// 	"io/ioutil"
// 	goLog "log"
// 	"sync"
// 	"testing"
//
// 	"github.com/corestoreio/errors"
// 	"github.com/corestoreio/log"
// 	"github.com/corestoreio/log/logw"
// 	"github.com/corestoreio/pkg/config"
// 	"github.com/corestoreio/pkg/store/scope"
// 	"github.com/stretchr/testify/assert"
// )
//
// // those tests cannot run in  because of reading and writing the debug log :-(
//
// var _ config.MessageReceiver = (*testSubscriber)(nil)
//
// type testSubscriber struct {
// 	t *testing.T
// 	f func(p config.Path) error
// }
//
// func (ts *testSubscriber) MessageConfig(p config.Path) error {
// 	//ts.t.Logf("Message: %s ScopeGroup %s ScopeID %d", p.String(), p.Scope.String(), p.ID)
// 	return ts.f(p)
// }
//
// func initLogger() (*log.MutexBuffer, log.Logger) {
// 	debugBuf := new(log.MutexBuffer)
// 	lg := logw.NewLog(
// 		logw.WithDebug(debugBuf, "testDebug: ", goLog.Lshortfile),
// 		logw.WithInfo(ioutil.Discard, "testInfo: ", goLog.Lshortfile),
// 	)
// 	lg.SetLevel(logw.LevelDebug)
// 	return debugBuf, lg
// }
//
// func TestPubSubBubbling(t *testing.T) {
//
// 	testPath := config.MustMakePath("aa/bb/cc")
//
// 	s := config.MustNewService(config.NewInMemoryStore(), config.WithPubSub())
//
// 	_, err := s.Subscribe(config.Route{}, nil)
// 	assert.True(t, errors.Empty.Match(err), "Error: %s", err)
//
// 	subID, err := s.Subscribe(testPath.Route, &testSubscriber{
// 		t: t,
// 		f: func(p config.Path) error {
// 			assert.Exactly(t, testPath.BindWebsite(123).String(), p.String(), "In closure Exactly")
// 			scp, id := p.ScopeID.Unpack()
// 			if scp == scope.Default {
// 				assert.Equal(t, int64(0), id)
// 			} else {
// 				assert.Equal(t, int64(123), id)
// 			}
// 			return nil
// 		},
// 	})
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")
//
// 	assert.NoError(t, s.Write(testPath.BindWebsite(123), 1))
// 	assert.NoError(t, s.Close())
//
// 	//t.Log("Before", "testPath", testPath.Route)
// 	testPath2 := testPath.Clone()
// 	assert.NoError(t, testPath2.Append(config.MakeRoute("Doh")))
// 	//t.Log("After", "testPath", testPath.Route, "testPath2", testPath2.Route)
//
// 	// send on closed channel
// 	assert.NoError(t, s.Write(testPath2.BindWebsite(3), 1))
// 	err = s.Close()
// 	assert.True(t, errors.AlreadyClosed.Match(err), "Error: %s", err)
// }
//
// func TestPubSubPanicSimple(t *testing.T) {
//
// 	debugBuf, logger := initLogger()
// 	s := config.MustNewService(config.NewInMemoryStore(), config.WithLogger(logger), config.WithPubSub())
// 	testPath := config.MakeRoute("xx/yy/zz")
//
// 	subID, err := s.Subscribe(testPath, &testSubscriber{
// 		t: t,
// 		f: func(_ config.Path) error {
// 			panic("Don't panic!")
// 		},
// 	})
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")
// 	assert.NoError(t, s.Write(config.MustMakePath(testPath).BindStore(123), 321), "Writing value 123 should not fail")
// 	assert.NoError(t, s.Close(), "Closing the service should not fail.")
// 	assert.Contains(t, debugBuf.String(), `config.pubSub.publish.recover.r recover: "Don't panic!"`)
// }
//
// func TestPubSubPanicError(t *testing.T) {
//
// 	debugBuf, logger := initLogger()
// 	s := config.MustNewService(config.NewInMemoryStore(), config.WithLogger(logger), config.WithPubSub())
//
// 	testPath := config.MustMakePath("aa/bb/cc")
//
// 	var pErr = errors.New("OMG! Panic!")
//
// 	subID, err := s.Subscribe(testPath, &testSubscriber{
// 		t: t,
// 		f: func(_ config.Path) error {
// 			panic(pErr)
// 		},
// 	})
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")
// 	assert.NoError(t, s.Write(config.MustMakePath(testPath).BindStore(123), 321))
//
// 	assert.NoError(t, s.Close())
// 	assert.Contains(t, debugBuf.String(), `config.pubSub.publish.recover.err error: "OMG! Panic!" route: "stores/123/aa/bb/cc"`)
// }
//
// func TestPubSubPanicMultiple(t *testing.T) {
//
// 	debugBuf, logger := initLogger()
// 	s := config.MustNewService(config.NewInMemoryStore(), config.WithLogger(logger), config.WithPubSub())
//
// 	subID, err := s.Subscribe(config.MustMakePath("xx"), &testSubscriber{
// 		t: t,
// 		f: func(p config.Path) error {
// 			assert.Equal(t, `xx/yy/zz`, p.Route.String())
// 			assert.Exactly(t, int64(987), p.ScopeID.ID())
// 			panic("One: Don't panic!")
// 		},
// 	})
// 	assert.NoError(t, err)
// 	assert.True(t, subID > 0)
//
// 	subID, err = s.Subscribe(config.MustMakePath("xx/yy"), &testSubscriber{
// 		t: t,
// 		f: func(p config.Path) error {
// 			assert.Equal(t, "xx/yy/zz", p.Route.String())
// 			assert.Exactly(t, int64(987), p.ScopeID.ID())
// 			panic("Two: Don't panic!")
// 		},
// 	})
// 	assert.NoError(t, err)
// 	assert.True(t, subID > 0)
//
// 	subID, err = s.Subscribe(config.MustMakePath("xx/yy/zz"), &testSubscriber{
// 		t: t,
// 		f: func(p config.Path) error {
// 			assert.Equal(t, "xx/yy/zz", p.Route.String())
// 			assert.Exactly(t, int64(987), p.ScopeID.ID())
// 			panic("Three: Don't panic!")
// 		},
// 	})
// 	assert.NoError(t, err)
// 	assert.True(t, subID > 0)
//
// 	assert.NoError(t, s.Write(config.MustMakePath("xx/yy/zz").BindStore(987), 789))
// 	assert.NoError(t, s.Close())
//
// 	assert.Contains(t, debugBuf.String(), `config.pubSub.publish.recover.r recover: "One: Don't panic!`)
// 	assert.Contains(t, debugBuf.String(), `config.pubSub.publish.recover.r recover: "Two: Don't panic!"`)
// 	assert.Contains(t, debugBuf.String(), `config.pubSub.publish.recover.r recover: "Three: Don't panic!"`)
// }
//
// func TestPubSubUnsubscribe(t *testing.T) {
//
// 	debugBuf, logger := initLogger()
// 	s := config.MustNewService(config.NewInMemoryStore(), config.WithLogger(logger), config.WithPubSub())
//
// 	var pErr = errors.New("WTF? Panic!")
// 	subID, err := s.Subscribe(config.MustMakePath("xx/yy/zz"), &testSubscriber{
// 		t: t,
// 		f: func(_ config.Path) error {
// 			panic(pErr)
// 		},
// 	})
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, subID, "The very first subscription ID should be 1")
// 	assert.NoError(t, s.Unsubscribe(subID))
// 	assert.NoError(t, s.Write(config.MustMakePath("xx/yy/zz").BindStore(123), 321))
// 	assert.NoError(t, s.Close())
// 	assert.Contains(t, debugBuf.String(), `config.Service.Write route: "stores/123/xx/yy/zz" val: 321`)
//
// }
//
// type levelCalls struct {
// 	sync.Mutex
// 	level2Calls int
// 	level3Calls int
// }
//
// func TestPubSubEvict(t *testing.T) {
//
// 	debugBuf, logger := initLogger()
// 	s := config.MustNewService(config.NewInMemoryStore(), config.WithPubSub(), config.WithLogger(logger))
//
// 	levelCall := new(levelCalls)
//
// 	var pErr = errors.New("WTF Eviction? Panic!")
//
// 	subID, err := s.Subscribe(config.MustMakePath("xx/yy"), &testSubscriber{
// 		t: t,
// 		f: func(p config.Path) error {
// 			assert.Contains(t, p.String(), "xx/yy")
// 			// this function gets called 3 times
// 			levelCall.Lock()
// 			levelCall.level2Calls++
// 			levelCall.Unlock()
// 			return nil
// 		},
// 	})
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, subID)
//
// 	subID, err = s.Subscribe(config.MustMakePath("xx/yy/zz"), &testSubscriber{
// 		t: t,
// 		f: func(p config.Path) error {
// 			assert.Contains(t, p.String(), "xx/yy/zz")
// 			levelCall.Lock()
// 			levelCall.level3Calls++
// 			levelCall.Unlock()
// 			// this function gets called 1 times and then gets removed
// 			panic(pErr)
// 		},
// 	})
// 	assert.NoError(t, err)
// 	assert.Equal(t, 2, subID)
//
// 	assert.NoError(t, s.Write(config.MustMakePath("xx/yy/zz").BindStore(123), 321))
// 	assert.NoError(t, s.Write(config.MustMakePath("xx/yy/aa").BindStore(123), 321))
// 	assert.NoError(t, s.Write(config.MustMakePath("xx/yy/zz").BindStore(123), 321))
//
// 	assert.NoError(t, s.Close())
//
// 	assert.Contains(t, debugBuf.String(), `config.pubSub.publish.recover.err error: "WTF Eviction? Panic!" route: "stores/123/xx/yy/zz"`)
//
// 	levelCall.Lock()
// 	assert.Equal(t, 3, levelCall.level2Calls)
// 	assert.Equal(t, 1, levelCall.level3Calls)
// 	levelCall.Unlock()
// 	err = s.Close()
// 	assert.True(t, errors.AlreadyClosed.Match(err), "Error: %s", err)
// }
