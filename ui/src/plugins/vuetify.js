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
import Vuetify from 'vuetify/lib';

import colors from 'vuetify/lib/util/colors';

import {
  VTooltip,
  VIcon,
  VListItem,
  VAvatar,
  VFlex,
  VLayout,
  VSelect,
  VTextarea,
  VCheckbox,
  VTextField,
  VFileInput
} from 'vuetify/lib'

Vue.component('v-tooltip', VTooltip);
Vue.component('v-icon', VIcon);
Vue.component('v-list-item', VListItem);
Vue.component('v-avatar', VAvatar);
Vue.component('v-flex', VFlex);
Vue.component('v-layout', VLayout);
Vue.component('v-select', VSelect);
Vue.component('v-text-field', VTextField);
Vue.component('v-textarea', VTextarea);
Vue.component('v-checkbox', VCheckbox);
Vue.component('v-file-input', VFileInput);

Vue.use(Vuetify);

export default new Vuetify({
  theme: {
    themes: {
      light: {
        primary: colors.red.lighten1,
        appbar: colors.grey.lighten5,
        sidebar: colors.blueGrey.darken3,
        secondary: colors.blue.lighten1,
        accent: colors.pink.darken2,
        background: '#fff', /*colors.blueGrey.lighten5,*/
        toolbar: '#000',
      },
      dark: {
        primary: colors.red.lighten1,
        secondary: colors.blue.lighten1,
        primaryText: colors.red.lighten1,
        toolbar: '#FFF',
        sidebar: colors.grey.darken4,
      },
    },
  },
});
