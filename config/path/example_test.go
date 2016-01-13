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

package path_test

//func Example() {
//
//	fmt.Println(path.MustNew("system/smtp/host").String())
//	fmt.Println(path.MustNew("system/smtp/host").Bind(scope.WebsiteID, 1).String())
//	// alternative way
//	fmt.Println(path.MustNew("system/smtp/host").BindStr(scope.StrWebsites, 1).String())
//
//	fmt.Println(path.MustNew("system/smtp/host").Bind(scope.StoreID, 3).String())
//	// alternative way
//	fmt.Println(path.MustNew("system/smtp/host").BindStr(scope.StrStores, 3).String())
//	// Group is not supported and falls back to default
//	fmt.Println(path.MustNew("system/smtp/host").Bind(scope.GroupID, 4).String())
//
//	// Output:
//	// default/0/system/smtp/host
//	// websites/1/system/smtp/host
//	// websites/1/system/smtp/host
//	// stores/3/system/smtp/host
//	// stores/3/system/smtp/host
//	// default/0/system/smtp/host
//}
