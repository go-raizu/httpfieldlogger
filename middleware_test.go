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

package httpfieldlogger_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-logr/logr/funcr"
	"github.com/stretchr/testify/assert"

	"github.com/go-raizu/httpfieldlogger"
)

func TestM(t *testing.T) {
	var handler http.Handler

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	l := funcr.New(func(prefix, args string) {
		fmt.Println(prefix, args)
	}, funcr.Options{})

	handler = httpfieldlogger.M(httpfieldlogger.New(l))(handler)

	assert.HTTPSuccess(t, handler.ServeHTTP, "GET", "", nil)
}
