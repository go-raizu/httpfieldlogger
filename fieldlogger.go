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
	"fmt"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/go-logr/logr"
	"github.com/go-raizu/harwp"
)

type FieldLogger struct {
	Logger logr.Logger
}

func (l *FieldLogger) NewLogEntry(r *http.Request) Event {
	fields := make([]any, 0, 10)

	fields = append(fields, "meth", r.Method)

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	fields = append(fields, "uri", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI))

	fields = append(fields, "proto", r.Proto)

	fields = append(fields, "from", r.RemoteAddr)

	entry := &Entry{Logger: l.Logger.WithValues(fields...)}
	return entry
}

func New(logger logr.Logger) *FieldLogger {
	return &FieldLogger{logger}
}

type Entry struct {
	Logger logr.Logger
}

func (e *Entry) WithField(k string, v any) Event {
	return &Entry{Logger: e.Logger.WithValues(k, v)}
}

func (e *Entry) Write(p harwp.ResponseWriterProxier, elapsed time.Duration) {
	e.Logger.WithValues(
		"status", p.StatusCode(),
		"length", humanize.IBytes(uint64(p.BytesWritten())),
		"elapsed", elapsed,
	).Info("")
}

func (e *Entry) Panic(v any) {
	var err error
	if err2, ok := v.(error); ok {
		err = err2
	}
	e.Logger.Error(err, "")
}
