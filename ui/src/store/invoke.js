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

import axios from 'axios';

export function invoke(method, url, arg, callback) {
  if (window.astilectron !== undefined) {
    window.astilectron.sendMessage({"name": url, "payload": arg, "method": method}, function (message) {
      console.log(message);
      if (message.payload !== undefined) {
        callback(message.payload);
      }
    });
  } else if (window.external.invoke !== undefined) {
    // local mode
    window.external.invoke(JSON.stringify(arg));
  } else {
    // server mode
    let server = '';
    if (typeof webpackHotUpdate !== 'undefined') {
      // dev mode
      server = `http://${window.location.hostname}:8081`;
    }
    axios({
      method,
      url: `${server}/api${url}`,
      data: arg,
    })
      .then((response) => {
        callback(response.data);
      });
  }
}
