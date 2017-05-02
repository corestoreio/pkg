// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dbr_test

//func TestSelectReturn(t *testing.T) {
//	ab := createRealSessionWithFixtures()
//
//	name, err := ab.Select("name").From("dbr_people").Where(Condition("email = 'jonathan@uservoice.com'")).ReturnString()
//	assert.NoError(t, err)
//	assert.Equal(t, name, "Jonathan")
//
//	count, err := ab.Select("COUNT(*)").From("dbr_people").ReturnInt64()
//	assert.NoError(t, err)
//	assert.Equal(t, count, int64(2))
//
//	names, err := ab.Select("name").From("dbr_people").Where(Condition("email = 'jonathan@uservoice.com'")).ReturnStrings()
//	assert.NoError(t, err)
//	assert.Equal(t, names, []string{"Jonathan"})
//
//	counts, err := ab.Select("COUNT(*)").From("dbr_people").ReturnInt64s()
//	assert.NoError(t, err)
//	assert.Equal(t, counts, []int64{2})
//}
