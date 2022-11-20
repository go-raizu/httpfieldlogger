// Copyright(c) 2022 individual contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// <https://www.apache.org/licenses/LICENSE-2.0>
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package httpfieldlogger

import (
	"net/http"
	"time"

	"github.com/go-raizu/ctxvalues"
	"github.com/go-raizu/harwp"
)

var logEntryCtxKey = ctxvalues.New[Event]()

type Event interface {
	WithField(k string, v any) Event
	Write(p harwp.ResponseWriterProxier, since time.Duration)
}

// GetEvent returns the in-context Event for a request.
func GetEvent(r *http.Request) Event {
	return logEntryCtxKey.GetOrZero(r.Context())
}

// WithEvent sets the in-context Event for a request.
func WithEvent(r *http.Request, entry Event) *http.Request {
	return r.WithContext(logEntryCtxKey.WithValue(r.Context(), entry))
}

type Logger interface {
	NewLogEntry(r *http.Request) Event
}

// M returns a logger handler middleware.
func M(l Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := harwp.NewResponseProxy(w)

			e := l.NewLogEntry(r)
			r = WithEvent(r, e)

			t1 := time.Now()
			next.ServeHTTP(p, r)
			d := time.Since(t1)

			e.Write(p, d)
		})
	}
}
