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

package ratelimit_test

//
//// Tests for a race condition
//func TestHTTPRateLimit_Concurrent_Map(t *testing.T) {
//
//	limiter, err := ratelimit.New()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	ctx := store.WithContextProvider(
//		context.Background(),
//		storemock.NewEurozzyService(
//			scope.MustSetByCode(scope.Website, "euro"),
//		),
//	)
//
//	a := mw.NewAdapter(
//		ctx,
//		mw.Chain(
//			ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
//				_, reqStore, err := store.FromContextProvider(ctx)
//				if err != nil {
//					return err
//				}
//				fmt.Fprintf(w, "hello world %s\nStore: %s", r.URL.RequestURI(), reqStore.Data.Code.String)
//				return nil
//			}),
//			storenet.WithInitStoreByFormCookie(), // extract the requested store code from the request
//			limiter.WithRateLimit(),              // then force the rate limit
//		),
//	)
//	ts := httptest.NewServer(a)
//
//	defer ts.Close()
//
//	var wg = sync.WaitGroup{}
//	for _, reqStore := range []string{"de", "at"} {
//		wg.Add(1)
//		go func(rsc string) {
//			defer wg.Done()
//			for i := 0; i < 6; i++ {
//				makeRequest(t, rsc, ts.URL)
//			}
//		}(reqStore)
//	}
//	wg.Wait()
//}
//
//func makeRequest(t *testing.T, requestedStore string, url string) {
//
//	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s=%s", url, storenet.HTTPRequestParamStore, requestedStore), nil)
//	if err != nil {
//		t.Fatal("NewRequest", err)
//	}
//	hc := &http.Client{
//		Timeout: time.Millisecond * 300,
//	}
//	res, err := hc.Do(req)
//	if err != nil {
//		t.Fatal("http.DefaultClient.Do", "err", err)
//	}
//
//	body, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		t.Fatal(err)
//	}
//	//t.Logf("Header: %#v", res.Header)
//	assert.Contains(t, string(body), fmt.Sprintf("hello world /?%s=%s\nStore: %s", storenet.HTTPRequestParamStore, requestedStore, requestedStore))
//	//t.Logf("Body: %#v", string(body))
//}
