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

package dbr

// IfNullAs returns IFNULL(`t1`.`c1`,`t2`.`c2`) AS `as` means
// if column c1 is null then use column c2.
func IfNullAs(t1, c1, t2, c2, as string) string {
	return NewAlias(
		"IFNULL("+Quoter.TableColumnAlias(t1, c1)[0]+", "+Quoter.TableColumnAlias(t2, c2)[0]+")",
		as,
	).String()
}
