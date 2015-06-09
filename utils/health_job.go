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

package utils

import "github.com/gocraft/health"

// HealthJobNoop default Job for the health.EventReceiver
type HealthJobNoop struct{}

var _ health.EventReceiver = (*HealthJobNoop)(nil)

func (HealthJobNoop) Event(eventName string)                                              {}
func (HealthJobNoop) EventKv(eventName string, kvs map[string]string)                     {}
func (HealthJobNoop) EventErr(eventName string, err error) error                          { return nil }
func (HealthJobNoop) EventErrKv(eventName string, err error, kvs map[string]string) error { return nil }
func (HealthJobNoop) Timing(eventName string, nanoseconds int64)                          {}
func (HealthJobNoop) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {}
