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
    <v-data-table
      :headers="headers"
      :items="logs"
      :loading="loading"
      :fixed-header="true"
      :footer-props="{'items-per-page-options': [25, 50, 100, 200]}"
      :items-per-page="100"
      dense
      :hide-default-footer="logs.length <= 100"
    >
      <template v-slot:body.prepend>
        <tr id="filter">
          <td v-for="h in headers"
              :key="h.text" role="columnheader"
              scope="col" style="">
            <v-text-field
              v-model="itemscol[h.value]"
              @keyup.enter.native="searchFilter"
              hide-details
              label="Filter"
              clearable
            />
          </td>
        </tr>
      </template>
      <template v-slot:item.time="{ item }">
        <div>{{ toLocal(item.time) }}</div>
      </template>
    </v-data-table>
</template>

<script>
  import {invoke} from "../store/invoke";
  import {DateTime} from "luxon";

  export default {
    name: 'logs',
    data() {
      return {
        timeFilter: '',
        fileFilter: '',
        messageFilter: '',
        loading: true,

        headers: [
          {
            text: "Time",
            value: "time",
            filter: f => {
              if (this.timeFilter === '' || this.timeFilter === null) return true
              else return (f + '').toLowerCase().includes(this['timeFilter'].toLowerCase())
            }
          },
          {
            text: "File",
            value: "file",
            filter: f => {
              if (this.fileFilter === '' || this.fileFilter === null) return true
              else return (f + '').toLowerCase().includes(this['fileFilter'].toLowerCase())
            }
          },
          {
            text: "Message",
            value: "message",
            filter: f => {
              if (this.messageFilter === '' || this.messageFilter === null) return true
              else return (f + '').toLowerCase().includes(this['messageFilter'].toLowerCase())
            }
          },
        ],
        logs: [],
        itemscol: {},
      };
    },

    computed: {
      checkFiltersEmpty() {
        return (this.timeFilter === '') && (this.fileFilter === '') && (this.messageFilter === '');
      },
    },

    methods: {
      loadFiles() {
        invoke('GET', '/logs', [], (data) => {
          this.logs = data.elements;
          this.loading = false;
        });
      },

      toLocal(s) {
        return DateTime.fromISO(s)
          .toLocaleString(DateTime.DATETIME_SHORT_WITH_SECONDS);
      },

      emptyFilter() {
        console.log('empty');
        this.itemscol = {}
        this.timeFilter = '';
        this.fileFilter = '';
        this.messageFilter = '';
        this.searchFilter();
      },

      searchFilter() {
        let column = '';

        for (const key in this.itemscol) {
          const value = this.itemscol[key];
          column = key;

          if (column === 'time') {
            this.timeFilter = value
          }

          if (column === 'file') {
            this.fileFilter = value
          }

          if (column === 'message') {
            this.messageFilter = value
          }
        }

      },

    },

    mounted() {
      this.loadFiles();
    },
  };
</script>
<style>
  .v-data-table {
    flex-grow: 1;
    overflow: hidden;
    display: grid;
    grid-template-rows: 1fr auto;
    height: 100%;
  }
  .v-data-table__wrapper {
    overflow: auto !important;
  }
</style>
