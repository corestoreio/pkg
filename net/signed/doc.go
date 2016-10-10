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

// Package signed provides a middleware to sign responses and validate requests.
//
// The middlewares can add the signature to the HTTP header or to the HTTP
// trailer or stores them internally. Store the hashes internally allows a
// transparent validation mechanism.
//
// With the use of HTTPS between two parties this package might not be needed.
//
// If you consider to use a 3rd untrusted participant then this package may help
// you. For example the 3rd participant is a mobile app. This app requests a
// product and a dynamic unique calculated price from the backend. The backend
// sends the product and its price to the mobile app. The app displays the
// information to the user. The user might add the product with that price to
// the cart by sending the product and the price back to the backend. The app
// uses the initial request which the app has received from the backend. The
// mobile app now simply forwards the unchanged first request bytes back to the
// backend. The backend can verify the request body by recalculating the hash
// found in the header.
//
// TODO(CyS) create a flowchart to demonstrate the usage.
//
// https://tools.ietf.org/html/draft-thomson-http-content-signature-00
// https://tools.ietf.org/html/draft-burke-content-signature-00
package signed
