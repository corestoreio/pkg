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

package mw

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/request"
	"github.com/rs/xstats"
	"github.com/zenazn/goji/web/mutil"
)

// Idea: github.com/rs/xaccess Copyright (c) 2015 Olivier Poitrey <rs@dailymotion.com> MIT License

// BlackholeXStat provides a type to disable the stats.
type BlackholeXStat struct{}

var _ xstats.XStater = (*BlackholeXStat)(nil)

// AddTag implements XStats interface
func (BlackholeXStat) AddTags(tags ...string) {}

// Gauge implements XStats interface
func (BlackholeXStat) Gauge(stat string, value float64, tags ...string) {}

// Count implements XStats interface
func (BlackholeXStat) Count(stat string, count float64, tags ...string) {}

// Histogram implements XStats interface
func (BlackholeXStat) Histogram(stat string, value float64, tags ...string) {}

// Timing implements xstats interface
func (BlackholeXStat) Timing(stat string, duration time.Duration, tags ...string) {}

// WithAccessLog is a middleware that logs all access requests performed on the
// sub handler and uses github.com/rs/xstats for collecting stats.
// Supported option arguments are: SetLogger().
// Provide none of those falls back to black hole logging. Log level must be
// set to info.
func WithAccessLog(x xstats.XStater, opts ...Option) Middleware {
	ob := newOptionBox(opts...)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Time request
			reqStart := time.Now()

			// Sniff the status and content size for logging
			lw := mutil.WrapWriter(w)

			h.ServeHTTP(lw, r)

			// Compute request duration
			reqDur := time.Since(reqStart)

			// Get request status
			status := ResponseStatus(r.Context(), lw.Status())

			// Log request stats

			tags := []string{
				"status:" + status,
				"status_code:" + strconv.Itoa(lw.Status()),
			}
			x.Timing("request_time", reqDur, tags...)
			x.Histogram("request_size", float64(lw.BytesWritten()), tags...)
			if ob.log.IsInfo() {
				ob.log.Info("request",
					log.String("proto", r.Proto),
					log.String("request_uri", r.RequestURI),
					log.String("method", r.Method),
					log.Stringer("uri", r.URL),
					log.String("type", "access"),
					log.String("status", status),
					log.Int("status_code", lw.Status()),
					log.Duration("duration", reqDur),
					log.String("requested-host", r.Host),
					log.Int("size", lw.BytesWritten()),
					log.Stringer("remote_addr", request.RealIP(r, request.IPForwardedTrust)),
					log.String("user_agent", r.Header.Get("User-Agent")),
					log.String("referer", r.Header.Get("Referer")),
				)
			}
		})
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
