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
  <div style="height: 100%;
    overflow: hidden;
    display: flex;
    flex-direction: column;">
    <v-breadcrumbs :items="breadcrumbs">
      <v-breadcrumbs-item
        slot="item"
        slot-scope="{ item }"
        exact
        @click="loadFiles( item.path )"
        class="breadcrumbsHover">
        {{ item.text }}
      </v-breadcrumbs-item>
    </v-breadcrumbs>
    <v-data-table
      :headers="headers"
      :items="files"
      :fixed-header="true"
      :footer-props="{'items-per-page-options': [10, 25, 50, 100]}"
      :items-per-page="25"
      :hide-default-footer="files.length <= 25"
      dense
    >
      <template v-slot:item.name="{ item }">
        <a v-if="item.dir" @click="loadFiles(item.path)">{{ item.name }}</a>
        <div v-else>{{ item.name }}</div>
      </template>
      <template v-slot:item.mtime="{ item }">
        <div>{{ toLocal(item.mtime) }}</div>
      </template>
      <template v-slot:item.size="{ item }">
        <div>{{ humanBytes(item.size, true) }}</div>
      </template>

      <template v-slot:item.actions="{ item }">
        <v-icon
          v-if="!item.dir"
          small
          @click="download(item)"
        >
          mdi-download
        </v-icon>
      </template>
    </v-data-table>
    <v-overlay
      :absolute="true"
      :opacity="0.5"
      :value="overlay"
      style="text-align: center"
    >
      <p>Saving file</p>
      <v-progress-circular indeterminate size="64"></v-progress-circular>
    </v-overlay>
  </div>
</template>

<script>
  import {invoke} from "../store/invoke";
  import {DateTime} from "luxon";

  export default {
    name: 'files',
    data() {
      return {
        overlay: false,
        breadcrumbs: [{text: '/', path: '/', link: true}],
        headers: [
          {text: "Name", value: "name"},
          // {text: "Path", value: "path"},
          {text: "Modified", value: "mtime"},
          {text: "Size", value: "size"},
          {text: "Download", value: "actions", sortable: false},
        ],
        files: []
      };
    },

    methods: {
      humanBytes(bytes, si) {
        const thresh = si ? 1000 : 1024;
        if (Math.abs(bytes) < thresh) {
          return `${bytes}B`;
        }
        const units = si
          ? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
          : ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
        let u = -1;
        do {
          bytes /= thresh;
          u += 1;
        } while (Math.abs(bytes) >= thresh && u < units.length - 1);
        return `${bytes.toFixed(1)}${units[u]}`;
      },

      toLocal(s) {
        return DateTime.fromISO(s)
          .toLocaleString(DateTime.DATETIME_SHORT_WITH_SECONDS);
      },

      loadFiles(path) {
        let breadcrumbs = [];
        let parts = path.split("/");
        for (let i = 0; i < parts.length; i += 1) {
          breadcrumbs.push({
            text: parts[i],
            href: parts[i],
            path: parts.slice(0, i + 1).join("/"),
            link: true,
          })
        }
        this.breadcrumbs = breadcrumbs;
        invoke('GET', '/files?path=' + path, [], (data) => {
          this.files = data.elements;
        });
      },

      download(item) {
        const that = this;
        if ('electron' in window) { // electron check
          window.electron.showSaveDialog({defaultPath: item.path.split('/').reverse()[0]}).then(result => {
            that.overlay = true;
            if (!result.canceled) {
              window.astilectron.sendMessage({
                "name": "save",
                "payload": {"src": item.path, "dest": result.filePath}
              }, function (message) {
                that.overlay = false;
              });
            }
          }).catch(err => {
            console.log(err)
          })
          return
        }

        this.$http.get("/api/file?path=" + item.path).then((response) => {
          const blob = new Blob([response.data])
          const link = document.createElement('a')
          link.href = URL.createObjectURL(blob)
          link.download = item.name
          link.click()
          URL.revokeObjectURL(link.href)
          console.log(response)
          that.data = response.data;
        }).catch(console.error);
      },
    },

    mounted() {
      this.loadFiles("/");
    },
  };
</script>
<style lang="scss">
  @import '~vuetify/src/styles/styles.sass';
  @import '../styles/colors.scss';
  @import '../styles/animations.scss';
  @import '~animate.css';

  .breadcrumbsHover:hover {
    color: $c-pink;
    cursor: pointer;
  }
</style>
