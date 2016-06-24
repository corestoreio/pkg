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

package backendratelimit_test

import "testing"

func TestHTTPRateLimit_WithConfig(t *testing.T) {

	//cfgStruct, err := ratelimit.NewConfigStructure()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//limiter, err := ratelimit.NewService(
	//	ratelimit.WithVaryBy(pathGetter{}),
	//	ratelimit.WithBackend(cfgStruct),
	//)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//cr := cfgmock.NewService(
	//	cfgmock.WithPV(cfgmock.PathValue{
	//		limiter.Backend.RateLimitBurst.MustFQ(scope.Website, 1):    0,
	//		limiter.Backend.RateLimitRequests.MustFQ(scope.Website, 1): 1,
	//		limiter.Backend.RateLimitDuration.MustFQ(scope.Website, 1): "i",
	//	}),
	//)
	//ctx := store.WithContextProvider(
	//	context.Background(),
	//	storemock.NewEurozzyService(
	//		scope.MustSetByCode(scope.Website, "euro"),
	//		store.WithStorageConfig(cr),
	//	),
	//)
	//
	//handler := limiter.WithRateLimit()(finalHandler200)
	//
	//runHTTPTestCases(t, ctx, handler, []httpTestCase{
	//	{"xx", 200, map[string]string{"X-Ratelimit-Limit": "1", "X-Ratelimit-Remaining": "0", "X-Ratelimit-Reset": "60"}},
	//	{"xx", 429, map[string]string{"X-Ratelimit-Limit": "1", "X-Ratelimit-Remaining": "0", "X-Ratelimit-Reset": "60", "Retry-After": "60"}},
	//})
}
