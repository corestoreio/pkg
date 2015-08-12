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

package user

import (
	"net/http"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/utils/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// JWTVerify is a middleware for echo to verify a JWT.
func JWTVerify(dbrSess dbr.SessionRunner) func(http.Handler) http.Handler {

	/*
		@todo
		1. load backend users from DB
		2. use them to check the valid token, etc
		3. create a polling service to update the cached backend user instead
		   of querying for each request the database.
		4. more stuff
	*/

	return func(c *echo.Context) error {
		// Skip WebSocket
		if (c.Request().Header.Get(echo.Upgrade)) == echo.WebSocket {
			return nil
		}
		token, err := jwt.ParseFromRequest(c.Request(), func(token *jwt.Token) (interface{}, error) {
			return []byte(`publicKey @todo`), nil
		})
		he := echo.NewHTTPError(http.StatusUnauthorized)

		if err != nil {
			log.Error("backend.JWTVerify.ParseFromRequest", "err", err, "req", c.Request())
			he.SetCode(http.StatusBadRequest)
			return he
		}

		if token.Valid {
			return nil
		}
		// log.Info() ?

		return he
	}
}
