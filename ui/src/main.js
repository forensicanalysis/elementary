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

import Vue from 'vue';
import axios from 'axios';
import VueAxios from 'vue-axios';
import VueLodash from 'vue-lodash';
import vuetify from './plugins/vuetify';


import App from './App.vue';
// import './registerServiceWorker';
import router from './router';
import store from './store';

import 'roboto-fontface/css/roboto/roboto-fontface.css';
import 'roboto-fontface/css/roboto-slab/roboto-slab-fontface.css';
import '@mdi/font/css/materialdesignicons.css';

Vue.use(VueAxios, axios);
Vue.use(VueLodash);

Vue.config.productionTip = false;

let v = new Vue({
  router,
  store,
  vuetify,
  render: h => h(App),
});
if (window && window.process && window.process.type) {
  document.addEventListener('astilectron-ready', function () {
    v.$mount('#app');
  });
} else {
  v.$mount('#app');
}
