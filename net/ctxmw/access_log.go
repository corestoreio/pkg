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

package ctxmw

import (
	"net/http"
	"strconv"
	"time"

	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/ctxlog"
	"github.com/corestoreio/csfw/net/httputil"
	"github.com/corestoreio/csfw/utils"
	"github.com/rs/xstats"
	"github.com/zenazn/goji/web/mutil"
	"golang.org/x/net/context"
)

// Idea: github.com/rs/xaccess Copyright (c) 2015 Olivier Poitrey <rs@dailymotion.com> MIT License

// WithAccessLog is a middleware that logs all access requests performed on the
// sub handler using github.com/corestoreio/csfw/net/ctxlog and
// github.com/rs/xstats stored in context.
func WithAccessLog() ctxhttp.Middleware {
	return func(h ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// Time request
			reqStart := time.Now()

			// Sniff the status and content size for logging
			lw := mutil.WrapWriter(w)

			err := h.ServeHTTPContext(ctx, lw, r)

			// Compute request duration
			reqDur := time.Since(reqStart)

			// Get request status
			status := ResponseStatus(ctx, lw.Status())

			// Log request stats
			sts := xstats.FromContext(ctx)
			tags := []string{
				"status:" + status,
				"status_code:" + strconv.Itoa(lw.Status()),
			}
			sts.Timing("request_time", reqDur, tags...)
			sts.Histogram("request_size", float64(lw.BytesWritten()), tags...)

			ctxlog.FromContext(ctx).Info("request",
				"error", utils.Errors(err),
				"method", r.Method,
				"uri", r.URL.String(),
				"type", "access",
				"status", status,
				"status_code", lw.Status(),
				"duration", reqDur.Seconds(),
				"size", lw.BytesWritten(),
				"remote_addr", httputil.GetRemoteAddr(r).String(),
				"user_agent", r.Header.Get("User-Agent"),
				"referer", r.Header.Get("Referer"),
			)
			return err
		}
	}
}

// ResponseStatus checks the context for timeout, canceled, ok or error.
func ResponseStatus(ctx context.Context, statusCode int) string {
	if ctx.Err() != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "timeout"
		}
		return "canceled"
	} else if statusCode >= 200 && statusCode < 300 {
		return "ok"
	}
	return "error"
}
