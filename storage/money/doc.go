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

/*
Package money uses a fixed-length guard for precision arithmetic.
Implements Un/Marshaller and Scan() method for database
columns including null, optimized for decimal(12, 4) fields.

Rounding is done on float64 to int64 by	the Rnd() function truncating
at values less than (.5 + (1 / Guardf))	or greater than -(.5 + (1 / Guardf))
in the case of negative numbers. The Guard adds four decimal places
of protection to rounding.
Decimal precision can be changed in the Precision() option
function. Precision() hold the places after the decimal place in the active
money struct field m.

http://en.wikipedia.org/wiki/Floating_point#Accuracy_problems

New()

Creating a new Currency struct:

	c := New()

Default values are 10000 for decimals, JSONLocale for Marshal/Unmarshal, Swedish
rounding is disabled and i18n.DefaultCurrency (en-US) for number and currency format.

The following options can be set while calling New():

	c := New(Swedish(Interval005), Guard(100), Precision(100))

Those values are really optional and even the order they appear ;-).
Default settings are:

	Precision 10000 which reflects decimal(12,4) database field
	Guard 	  10000 which reflects decimal(12,4) database field
	Swedish   No rounding

If you need to temporarily set a different option value you can stick to this pattern:
http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html

	prev := m.Option(Swedish(Interval005))
	defer m.Option(prev)
	// do something with the different Swedish rounding

Initial Idea: Copyright (c) 2011 Jad Dittmar
https://github.com/Confunctionist/finance

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/
package money
