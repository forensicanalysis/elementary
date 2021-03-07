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

package server

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/forensicanalysis/forensicstore"
)

const (
	cacheExpiration      = 5 * time.Minute
	cacheCleanupInterval = 10 * time.Minute
)

var queryCache *cache.Cache // nolint: gochecknoglobals

func init() { // nolint: gochecknoinits
	queryCache = cache.New(cacheExpiration, cacheCleanupInterval)
}

func storequery(store *forensicstore.ForensicStore, q string) ([]forensicstore.JSONElement, error) {
	elems, found := queryCache.Get(q)
	if !found {
		fmt.Println(q)
		elems, err := store.Query(q)
		if err != nil {
			return nil, err
		}
		queryCache.Set(q, elems, cache.DefaultExpiration)
		return elems, nil
	}

	return elems.([]forensicstore.JSONElement), nil
}
