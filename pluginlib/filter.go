// Copyright (c) 2020 Siemens AG
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// Author(s): Jonas Plum

package pluginlib

import (
	"strings"

	"github.com/tidwall/gjson"
)

// A Filter is a list of mappings that should be used for a Task.
type Filter []map[string]string

// Match tests if an element matches the filter.
func (f Filter) Match(element []byte) bool {
	if len(f) == 0 {
		return true
	}
	for _, condition := range f {
		if f.matchCondition(condition, element) {
			return true
		}
	}
	return false
}

func (f Filter) matchCondition(condition map[string]string, element []byte) bool {
	for attribute, value := range condition {
		if !strings.Contains(gjson.GetBytes(element, attribute).String(), value) {
			return false
		}
	}
	return true
}

func ExtractFilter(filtersets []string) Filter {
	filter := Filter{}
	for _, filterset := range filtersets {
		filterelement := map[string]string{}
		for _, kv := range strings.Split(filterset, ",") {
			kvl := strings.SplitN(kv, "=", 2)
			if len(kvl) == 2 {
				filterelement[kvl[0]] = kvl[1]
			}
		}

		filter = append(filter, filterelement)
	}
	return filter
}
