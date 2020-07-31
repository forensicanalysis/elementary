<!--
Copyright (c) 2020 Siemens AG

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

Author(s): Jonas Plum
-->
<template>
  <div>
    <input type="text" v-model="path" style="width: 100%;">
    <v-tabs v-model="tab" icons-and-text>
      <v-tabs-slider/>
      <v-tab v-if="_.endsWith(path, '.pdf')" href="#tab-pdf">PDF
        <v-icon v-text="'mdi-file-pdf-outline'"/>
      </v-tab>
      <v-tab v-if="isImagePath(path)" href="#tab-image">Image
        <v-icon v-text="'mdi-file-image-outline'"/>
      </v-tab>
      <v-tab v-if="isTextPath(path)" href="#tab-text">Text
        <v-icon v-text="'mdi-file-document-outline'"/>
      </v-tab>
      <v-tab href="#tab-hex">Hex
        <v-icon v-text="'mdi-file-document-outline'"/>
      </v-tab>
    </v-tabs>
    <v-tabs-items v-model="tab">
      <v-tab-item :value="'tab-text'">
        <v-card flat class="pa-2" style="overflow: auto">
          <pre>{{data}}</pre>
        </v-card>
      </v-tab-item>
      <v-tab-item :value="'tab-hex'" style="overflow: auto">
        <v-card flat class="pa-2">{{ hex }}</v-card>
      </v-tab-item>
      <v-tab-item :value="'tab-pdf'" style="overflow: auto">
        <pdf
          ref="pdf"
          :src="data"
          :page="1"
          style="border: 1px solid red"
          @num-pages="numPages = $event"
          @loaded="loaded = true"
        />
      </v-tab-item>
      <v-tab-item :value="'tab-image'" style="overflow: auto">
        <v-card flat class="pa-2" style="overflow: auto; max-height: 100%">
          <img
            :src="'http://localhost:8081/api/file?path='+path"
            style="width: 100%" alt=""/>
        </v-card>
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>

<script>
import pdf from 'vue-pdf';
// import { component as VueCodeHighlight } from 'vue-code-highlight';
// import JsonToHtml from '@/components/json-to-html';
// import hexyjs from 'hexyjs';

export default {
  name: 'file',
  components: {
    // JsonToHtml,
    pdf,
    // VueCodeHighlight,
  },
  data() {
    return {
      path: '',
      tab: null,
      data: null,
      hex: '0000',
      active: 'Info',
      loaded: false,
    };
  },
  watch: {
    path() {
      const vm = this;
      this.$http.get(`http://localhost:8081/api/file?path=${this.path}`)
        .then((response) => {
          // vm.hex = hexyjs.strToHex(response.data);
          vm.data = response.data;
        });
    },
  },
  props: {
    /* future_path: {
        type: Object,
        required: true,
      }, */
  },
  mounted() {
    this.path = 'WindowsRecycleBin/$I9F219A.jpg';
  },
  methods: {
    hasEnding(path, endings) {
      for (let i = 0; i < endings.length; i += 1) {
        if (Vue._.endsWith(path, `.${endings[i]}`)) {
          return true;
        }
      }
      return false;
    },
    isImagePath(path) {
      return this.hasEnding(path, ['jpg', 'jpeg', 'png', 'svg']);
    },
    isTextPath(path) {
      return this.hasEnding(path, ['txt', 'md']);
    },
  },
};
</script>
