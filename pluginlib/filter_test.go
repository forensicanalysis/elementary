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
	"testing"
)

func TestFilter_Match(t *testing.T) {
	type args struct {
		element string
	}
	tests := []struct {
		name string
		f    Filter
		args args
		want bool
	}{
		{"simple match", Filter{{"name": "foo"}}, args{`{"name": "foo"}`}, true},
		{"no match", Filter{{"name": "foo"}}, args{`{"name": "bar"}`}, false},
		{"nil filter", nil, args{`{"name": "foo"}`}, true},
		{"contains match", Filter{{"name": "foo"}}, args{`{"name": "xfool"}`}, true},
		{"simple match", Filter{{"name": "foo"}}, args{`{"name": "foo", "bar": "baz"}`}, true},
		{"multi match", Filter{{"name": "foo", "bar": "baz"}}, args{`{"name": "foo", "bar": "baz"}`}, true},
		{"any match", Filter{{"x": "y"}, {"name": "foo", "bar": "baz"}}, args{`{"name": "foo", "bar": "baz"}`}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Match([]byte(tt.args.element)); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
