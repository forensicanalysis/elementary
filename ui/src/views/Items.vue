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
  <div class="d-flex flex-row" style="overflow: hidden; height: 100%;">
    <div
      ref="drawerLeft"
      class="flex-grow-0 flex-shrink-0 verticalbar scrollableArea"
      :style="'width: ' + leftWidth + 'px'"
      style="border-right: 1px solid rgba(0, 0, 0, 0.12); overflow-x: hidden !important;"
    >
      <v-expansion-panels v-model="panel" accordion multiple>
        <v-expansion-panel>
          <v-expansion-panel-header>
            <v-subheader class="navigationHeader tr-2">
              <transition name="fade-fast">
                <a v-show="leftExtended" class="pl-4">TYPE</a>
              </transition>
              <v-spacer></v-spacer>
              <v-subheader class="tr-2">
                <v-icon :class="{'navigationMenuIcon': !leftExtended}">mdi-format-list-bulleted-type
                </v-icon>
              </v-subheader>
            </v-subheader>
          </v-expansion-panel-header>
          <v-expansion-panel-content>
            <v-list style="padding: 0 !important" dense>
              <hr class="divider"/>
              <div v-if="loadingType">
                <v-spacer/>
                <v-progress-linear
                  indeterminate
                  color="#EF5350"
                  :height="2"
                ></v-progress-linear>
              </div>
              <v-list-item-group v-else v-model="itemTypeIndex" color="primary">
                <v-list-item @click="loadType('')">
                  <v-list-item-icon>
                    <v-icon>mdi-file-multiple</v-icon>
                  </v-list-item-icon>
                  <v-list-item-content>
                    <v-list-item-title>All</v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
                <v-list-item
                  v-for="(table, i) in _.sortBy(tables, ['name'])"
                  :key="i"
                  @click="loadType(table['name'])"
                  class="px-4"
                >
                  <v-list-item-icon>
                    <v-icon
                      v-if="_.has(templates[table['name']], 'icon')"
                      v-text="'mdi-'+templates[table['name']].icon"/>
                    <v-icon v-else>mdi-file-outline</v-icon>
                  </v-list-item-icon>
                  <v-list-item-content>
                    <v-list-item-title v-text="_.startCase(table.name)"/>
                  </v-list-item-content>
                  <v-list-item-icon v-show="leftExtended">
                    <span v-if="'count' in table">{{table['count']}}</span>
                  </v-list-item-icon>
                </v-list-item>
              </v-list-item-group>
            </v-list>
          </v-expansion-panel-content>
        </v-expansion-panel>
        <v-expansion-panel v-show="showLocation">
          <v-expansion-panel-header>
            <v-subheader class="navigationHeader tr-2">
              <transition name="fade-fast">
                <a v-show="leftExtended" class="pl-4">LOCATION</a>
              </transition>
              <v-spacer></v-spacer>
              <v-subheader class="tr-2">
                <v-icon :class="{'navigationMenuIcon': !leftExtended}">mdi-file-tree</v-icon>
              </v-subheader>
            </v-subheader>
          </v-expansion-panel-header>
          <v-expansion-panel-content>
            <v-list style="padding: 0 !important" dense>
              <v-list-item-group>
                <hr class="divider"/>
                <div v-if="loadingTree">
                  <v-spacer/>
                  <v-progress-linear
                    indeterminate
                    color="#EF5350"
                    :height="2"
                  ></v-progress-linear>
                </div>
                <div v-else>
                  <transition name="fade-slow" mode="out-in">
                    <div v-if="!leftExtended" key="dots">
                      <v-icon class="tr-2"
                              style="margin-left: 16px">
                        mdi-dots-horizontal
                      </v-icon>
                    </div>
                    <div v-else key="tree">
                      <v-treeview
                        class="tr-2"
                        v-show="directories.length > 0 && leftExtended"
                        activatable
                        hoverable
                        dense
                        transition
                        @update:active="updatedir"
                        :active.sync="active"
                        :items="directories"
                        :load-children="fetchTreeChildren"
                        :open.sync="open"
                        item-key="path"
                        color="primary"
                        ref="treeView"
                      >
                        <template v-slot:prepend="{ item, open }">
                          <v-icon v-if="item.children">
                            {{ open ? 'mdi-folder-open' : 'mdi-folder' }}
                          </v-icon>
                          <v-icon v-if="!item.children">mdi-folder-outline</v-icon>
                        </template>
                      </v-treeview>
                    </div>
                  </transition>
                </div>
              </v-list-item-group>
            </v-list>
          </v-expansion-panel-content>
        </v-expansion-panel>
        <!--<v-expansion-panel>
          <v-expansion-panel-header>
            <v-subheader class="navigationHeader tr-2">
              <transition name="fade-fast">
                <a v-show="leftExtended" class="pl-4">LABELS</a>
              </transition>
              <v-spacer></v-spacer>
              <v-subheader class="tr-2">
                <v-icon :class="{'navigationMenuIcon': !leftExtended}">mdi-label-outline</v-icon>
              </v-subheader>
            </v-subheader>
          </v-expansion-panel-header>
          <v-expansion-panel-content>
            <v-list style="padding: 0 !important" dense>
              <v-list-item-group>
                <hr class="divider"/>
                <div v-if="loadingLabels">
                  <v-spacer/>
                  <v-progress-linear
                    indeterminate
                    color="#EF5350"
                    :height="2"
                  ></v-progress-linear>
                </div>
                <div v-else style="margin-top: 10px; margin-left: 10px">
                  <transition name="fade-slow" mode="out-in">
                    <div v-if="!leftExtended" key="dots">
                      <v-icon class="tr-2"
                              style="margin-left: 6px">
                        mdi-dots-horizontal
                      </v-icon>
                    </div>
                    <div v-if="!leftExtended">
                    </div>
                    <div v-else key="labels">
                      <v-chip-group
                        v-model="selectedLabel"
                        column
                        active-class="primary--text"
                      >
                        <v-chip
                          class="ma-2 lighten-3 navigationChip"
                          color="grey"
                          text-color="#424242"
                          v-for="label in labels"
                          :key="label"
                          :value="label"
                          small
                          style="margin: 2px !important"
                          @click="filterLabels(label)"
                        >
                          <v-avatar
                            left
                            class="grey lighten-2"
                            style="margin-left: -12px !important; color: #EF5350"
                          >
                            {{countOccurrences(labels, label)}}
                          </v-avatar>
                          {{label}}
                        </v-chip>
                      </v-chip-group>
                    </div>
                  </transition>
                </div>
              </v-list-item-group>
            </v-list>
          </v-expansion-panel-content>
        </v-expansion-panel>-->
      </v-expansion-panels>
      <v-divider/>
      <div class="d-flex justify-end">
        <v-btn small icon @click="toogleLeft" class="menu-hide">
          <v-icon class="navigationDrawerIcon"
                  :class="{'rotate180': !leftExtended}">mdi-chevron-left
          </v-icon>
        </v-btn>
      </div>
    </div>
    <div class="flex-grow-1 d-flex flex-row" style="width: 100%; overflow: hidden">
      <div style="overflow: hidden; transition: width 0.2s ease-in; display: flex; flex-direction: column;"
           class="pt-3 flex-shrink-1 flex-grow-1"
      >
        <v-bottom-sheet v-model="bottomSheet"
                        inset
                        hide-overlay
                        no-click-animation
                        persistent>
          <v-sheet class="text-center bottomSheetStyle"
                   height="auto"
                   style="background: #EF5350">
            <v-btn
              small
              icon
              @click="labelsDialogMultiple = true"
            >
              <v-icon class="bottomSheetIcon">mdi-label-outline</v-icon>
            </v-btn>
          </v-sheet>
        </v-bottom-sheet>

        <v-text-field
          v-model="search"
          prepend-inner-icon="mdi-magnify"
          clear-icon="mdi-close-circle"
          clearable
          dense
          outlined
          @keyup.enter.native="searchFilter"
          @click:append="searchFilter"
          @click:clear="clearFilter"
          class="mx-3"
        />
        <v-data-table
          v-model="selectedItems"
          :headers="headers"
          :items="items"
          :loading="loadElementsTable"
          :server-items-length="itemCount"
          @update:options="updateopt"
          :fixed-header="true"
          @click:row="select"
          :footer-props="{'items-per-page-options': [25, 50, 100, 200, 500]}"
          :items-per-page="100"
          :hide-default-footer="itemCount <= 100"
          show-select
          dense
        >
          <template v-slot:body.prepend>
            <tr>
              <td @click="emptyFilter">
                <v-icon v-if="_.isEmpty(filter.columns) && filteredLabel === ''" color="primary"
                        small class="ml-1">
                  mdi-filter-outline
                </v-icon>
                <v-icon v-else color="primary" small class="ml-1">
                  mdi-filter-remove-outline
                </v-icon>
              </td>
              <td v-for="h in headers"
                  :key="h.text" role="columnheader"
                  scope="col">
                <v-text-field
                  v-model="itemscol[h.value]"
                  @keyup.enter.native="searchFilter"
                  hide-details
                  label="Filter"
                  :reverse="h.align === 'right'"
                  clearable
                />
              </td>
            </tr>
          </template>

          <template v-slot:item.title="{ item }">
            <div>{{ title(item) }}</div>
          </template>
          <template v-slot:item.atime="{ item }">
            <div>{{ toLocal(item.atime) }}</div>
          </template>
          <template v-slot:item.mtime="{ item }">
            <div>{{ toLocal(item.mtime) }}</div>
          </template>
          <template v-slot:item.ctime="{ item }">
            <div>{{ toLocal(item.ctime) }}</div>
          </template>
          <template v-slot:item.origin="{ item }">
            <div><span v-if="'path' in item.origin">{{item.origin.path}}</span></div>
          </template>
          <template v-slot:item.size="{ item }">
            <div>{{ humanBytes(item.size, true) }}</div>
          </template>
        </v-data-table>
      </div>
      <div
        :style="'width: ' + rightWidth + '%'"
        style="border-left: 1px solid rgba(0, 0, 0, 0.12)"
        class="flex-grow-0 flex-shrink-0 verticalbar scrollableArea"
        ref="drawerRight"
      >
        <v-dialog
          v-model="labelsDialogMultiple"
          max-width="800px"
          height="500px"
          style="overflow: hidden"
        >
          <v-card>
            <v-card-title class="labelsDialogTitle">
              LABELS
            </v-card-title>
            <hr class="divider" style="margin: 0 !important"/>
            <v-combobox
              v-model="labelsToAddMultiple"
              :items="labels"
              :search-input.sync="searchLabelsMultiple"
              hide-selected
              label="Enter new labels or select existing ones."
              multiple
              chips
              outlined
              style="margin: 20px"
              dense
              deletable-chips
            >
              <template v-slot:no-data>
                <v-list-item>
                  <v-list-item-content>
                    <v-list-item-title>
                      No results matching. Press <kbd>enter</kbd> to
                      create a new one.
                    </v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </template>
            </v-combobox>
            <v-card-actions>
              <v-btn
                color="primary"
                text
                @click="discardAddLabels(true)"
              >
                Close
              </v-btn>
              <v-spacer></v-spacer>
              <v-btn
                color="primary"
                text
                @click="labelsConfirmationMultiple = !labelsConfirmationMultiple"
              >
                Confirm
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
        <v-dialog
          v-model="labelsConfirmationMultiple"
          max-width="500px"
        >
          <v-card>
            <v-card-title class="labelsDialogTitle">
              <span>Are You Sure?</span>
              <v-spacer></v-spacer>
            </v-card-title>
            <hr class="divider" style="margin: 0 !important"/>
            <div style="padding: 12px">
              <p style="margin: 4px" class="noLabelsText">Following labels will be added to
                {{this.selectedItems.length}}
                items:</p>
              <v-chip style="margin: 4px" v-for="c in labelsToAddMultiple">{{c}}</v-chip>
            </div>
            <hr class="divider" style="margin: 0 !important"/>
            <v-card-actions>
              <v-btn
                color="primary"
                text
                @click="labelsConfirmationMultiple = false"
              >
                No
              </v-btn>
              <v-spacer/>
              <v-btn
                color="primary"
                text
                @click="setLabelsMultipleChunk(selectedItems, labelsToAddMultiple)"
              >
                Yes
              </v-btn>
            </v-card-actions>
            <v-progress-linear
              v-if="loadingLabelsOperation"
              color="#EF5350"
              buffer-value="0"
              stream
            ></v-progress-linear>
          </v-card>
        </v-dialog>
        <v-dialog
          v-model="labelsDialog"
          max-width="800px"
          height="500px"
        >
          <v-card>
            <v-card-title class="labelsDialogTitle">
              LABELS
            </v-card-title>
            <hr class="divider" style="margin: 0 !important"/>
            <div
              style="padding: 2px 14px !important; max-height: 40px !important; min-height: 40px !important; background: whitesmoke">
              <v-chip-group
                v-model="labelsToRemove"
                multiple
                show-arrows
              >
                <v-chip class="lighten-3"
                        color="grey"
                        text-color="#EF5350"
                        small
                        filter
                        filter-icon="mdi-close"
                        style="margin: 2px !important"
                        v-for="label in existingLabels" :key="label">
                  {{ label }}
                </v-chip>
                <p class="noLabelsText"
                   v-if="existingLabels.length === 0 && labelsToAdd.length === 0">No labels
                  exist.</p>
                <v-chip class="lighten-4"
                        color="green"
                        text-color="#EF5350"
                        small
                        filter
                        filter-icon="mdi-close"
                        style="margin: 2px !important"
                        v-for="label in labelsToAdd" :key="label">
                  {{ label }}
                </v-chip>
              </v-chip-group>
            </div>
            <hr class="divider" style="margin: 0 !important"/>
            <v-combobox
              v-model="labelsToAdd"
              :items="removeExistingLabels(existingLabels)"
              :search-input.sync="searchLabels"
              hide-selected
              label="Enter new labels or select existing ones."
              multiple
              chips
              outlined
              style="margin: 20px"
              dense
              deletable-chips
            >
              <template v-slot:no-data>
                <v-list-item>
                  <v-list-item-content>
                    <v-list-item-title>
                      No results matching. Press <kbd>enter</kbd> to
                      create a new one.
                    </v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </template>
            </v-combobox>
            <v-card-actions>
              <v-btn
                color="primary"
                text
                @click="discardAddLabels(false)"
              >
                Close
              </v-btn>
              <v-spacer></v-spacer>
              <v-btn
                color="primary"
                text
                @click="labelsConfirmation = !labelsConfirmation"
              >
                Confirm
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
        <v-dialog
          v-model="labelsConfirmation"
          max-width="500px"
        >
          <v-card v-if="this.labelsToRemoveElements">
            <v-card-title class="labelsDialogTitle">
              <span>Are You Sure?</span>
              <v-spacer></v-spacer>
            </v-card-title>

            <div v-if="labelsToAdd.length !== 0">
              <hr class="divider" style="margin: 0 !important"/>
              <div style="padding: 12px">
                <p style="margin: 4px" class="noLabelsText">Following labels will be added:</p>
                <v-chip style="margin: 4px" v-for="c in labelsToAdd">{{c}}</v-chip>
              </div>
            </div>

            <div v-if="calculateLabelsToRemove.length !== 0">
              <hr class="divider" style="margin: 0 !important"/>
              <div style="padding: 12px">
                <p style="margin: 4px" class="noLabelsText">Following labels will be removed:</p>
                <v-chip style="margin: 4px" v-for="c in calculateLabelsToRemove">{{c}}</v-chip>
              </div>
            </div>

            <div v-if="labelsToAdd.length === 0 && calculateLabelsToRemove.length === 0">
              <hr class="divider" style="margin: 0 !important"/>
              <div style="padding: 12px">
                <p style="margin: 4px" class="noLabelsText">No changes made.</p>
              </div>
            </div>

            <hr class="divider" style="margin: 0 !important"/>
            <v-card-actions>
              <v-btn
                color="primary"
                text
                @click="labelsConfirmation = false"
              >
                No
              </v-btn>
              <v-spacer/>
              <v-btn
                color="primary"
                text
                @click="setLabelsMultiple(item.id, labelsToAdd, calculateLabelsToRemove)"
              >
                Yes
              </v-btn>
            </v-card-actions>
            <v-progress-linear
              v-if="loadingLabelsOperation"
              color="#EF5350"
              buffer-value="0"
              stream
            ></v-progress-linear>
          </v-card>
        </v-dialog>
        <v-toolbar
          v-if="!_.isEmpty(item)"
          class="elevation-0 ml-2 mr-0"
          dense>
          <v-toolbar-title>{{ item.id }}</v-toolbar-title>
          <v-spacer/>
          <v-toolbar-items>
            <!--v-btn
              small
              icon
              @click="labelsDialog = true"
            >
              <v-icon class="detailsIcon">mdi-label-outline</v-icon>
            </v-btn-->
            <v-btn
              small
              icon
              @click="expandDetails">
              <v-icon v-if="this.rightWidth === 50" class="detailsIcon">mdi-fullscreen</v-icon>
              <v-icon v-if="this.rightWidth === 100" class="detailsIcon">mdi-fullscreen-exit
              </v-icon>
            </v-btn>
            <v-btn
              small
              icon
              @click="rightNull"
            >
              <v-icon class="detailsIcon">mdi-close</v-icon>
            </v-btn>
          </v-toolbar-items>
        </v-toolbar>
        <v-divider class="mx-0"/>
        <item :content="item" :labels="labels"/>
      </div>
    </div>
  </div>
</template>

<script>
  import {DateTime} from 'luxon';
  import item from '@/views/Document.vue';
  import {invoke} from "../store/invoke";
  import {templates} from "../store/templates";
  import Vue from "vue";

  export default {

    name: 'items',

    components: {
      item,
    },

    data() {
      return {
        type: '',
        items: [],
        itemCount: 0,
        headers: [],
        sort: {
          type: '',
          columns: {},
        },
        filter: {
          type: '',
          columns: {},
        },
        limit: 25,
        offset: 0,
        templates,

        item: {},
        listPane: 80,

        panel: [0, 1, 2],

        leftExtended: true,
        itemscol: {},
        rightWidth: 0,
        leftWidth: 250,

        loadingType: true,
        loadingTree: true,
        loadingLabels: true,
        loadElementsTable: true,
        loadingLabelsOperation: true,

        labelsDialog: false,
        labelsConfirmation: false,
        labelsToAdd: [],
        labelsToRemove: [],
        labelsToRemoveElements: [],
        searchLabels: '',

        labelsDialogMultiple: false,
        labelsConfirmationMultiple: false,
        labelsToAddMultiple: [],
        searchLabelsMultiple: '',

        selectedItems: [],
        label: '',
        itemID: '',

        selectedLabel: null,
        filteredLabel: '',

        search: '',
        active: [],
        open: [],

        itemTypeIndex: 0,
        // options: {},

        directories: [],
        labels: [],
        tables: [],
      };
    },

    computed: {
      count() {
        for (let i = 0; i < this.tables.length; i += 1) {
          if (this.type === this.tables[i].name) {
            return this.tables[i].count;
          }
        }
        return 0;
      },

      showLocation() {
        const type = this.type;
        return (type === 'file' || type === 'directory' || type === 'windows-registry-key')
      },

      calculateLabelsToRemove() {
        if (this.item.label) {
          const elements = this.existingLabels.concat(this.labelsToAdd);
          const toRemove = [];
          for (let i = 0; i < this.labelsToRemove.length; i++) {
            toRemove.push(elements[this.labelsToRemove[i]]);
          }
          return toRemove
        } else {
          const elements = this.labelsToAdd;
          const toRemove = [];
          for (let i = 0; i < this.labelsToRemove.length; i++) {
            toRemove.push(elements[this.labelsToRemove[i]]);
          }
          return toRemove
        }
      },

      existingLabels() {
        if (this.item.label) {
          return Object.keys(this.removeFalseLabels(this.item.label));
        } else {
          return [];
        }
      },

      bottomSheet() {
        return this.selectedItems.length !== 0;
      },
    },
    methods: {
      toogleLeft() {
        if (this.leftWidth !== 56) {
          this.leftWidth = 56;
          this.leftExtended = false;
        } else {
          this.leftWidth = 250;
          this.leftExtended = true;
        }
      },

      rightNull() {
        this.rightWidth = 0;
      },
      rightHalf() {
        this.rightWidth = 50;
      },
      rightFull() {
        this.rightWidth = 100;
      },

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

      title(element) {
        if (_.has(this.templates, element['type'])) {
          return element[this.templates[element['type']].headers[0].value];
        }
        if (_.has(element, 'name')) {
          return element['name'];
        }
        if (_.has(element, 'key')) {
          return element['key'];
        }
        if (_.has(element, 'title')) {
          return element['title'];
        }
        return '';
      },

      getDirectories() {
        this.loadingTree = true;
        let that = this;
        if (this.directories.length === 0) {
          this.loadDirectories('', (directories) => {
            that.directories = directories
          });
          this.loadingTree = false;
        }
      },

      loadType(table) {
        this.loadingTree = true;

        this.directories = [];
        this.itemscol = {};
        this.active = [];
        this.open = [];

        this.item = {};
        this.rightNull();
        this.sort = {
          table: '',
          columns: {},
        };
        this.filter = {
          type: '',
          columns: {},
        };
        this.type = table;
        console.log("li1");
        this.loadItems();
        for (let i = 0; i < this.headers.length; i += 1) {
          this.itemscol[this.headers[i].value] = '';
        }

        this.getDirectories();

        this.filter = {
          type: this.type,
          columns: {},
        };
      },

      emptyFilter() {
        this.itemscol = {}
        this.filter = {
          type: this.type,
          columns: {},
        };
        this.filteredLabel = '';
        this.label = '';
        this.searchFilter()
      },

      select(e) {
        this.item = e;
        this.rightHalf();
      },

      toLocal(s) {
        return DateTime.fromISO(s)
          .toLocaleString(DateTime.DATETIME_SHORT_WITH_SECONDS);
      },

      updateopt(e) {
        this.offset = (e.page - 1) * e.itemsPerPage;
        this.limit = e.itemsPerPage;
        const sort = {
          table: this.type,
          columns: {},
        };
        for (let i = 0; i < e.sortBy.length; i += 1) {
          if (e.sortDesc[i]) {
            sort.columns[e.sortBy[i]] = 'DESC';
          } else {
            sort.columns[e.sortBy[i]] = 'ASC';
          }
        }
        this.sort = sort;
        console.log("li2");
        this.loadItems();
      },

      updatedir(e) {
        console.log("updatedir");
        let slash = '';
        let column = '';

        if (this.type === 'directory') {
          slash = '/';
          column = 'path';
        } else if (this.type === 'file') {
          slash = '/';
          column = 'origin.path';
        } else if (this.type === 'windows-registry-key') {
          slash = '\\';
          column = 'key';
        } else {
          console.log('TABLE DOES NOT EXIST!');
        }

        this.filter.type = this.type;
        if (e.length === 0) {
          this.filter.columns[column] = [];
        } else {
          this.filter.columns[column] = [e[0] + slash];
        }

        console.log("li3");
        this.loadItems();
      },

      searchFilter() {
        this.filter.type = this.type;

        let column = '';

        for (const key in this.itemscol) {
          const value = this.itemscol[key];
          if (key === 'origin') {
            column = 'origin.path';
          } else {
            column = key;
          }
          this.filter.columns[column] = [value];
        }

        if (this.search !== '') {
          this.filter.columns.elements = [this.search];
        }
        console.log("li4");
        this.loadItems();
      },

      clearFilter() {
        this.search = '';
        this.searchFilter();
      },

      expandDetails() {
        if (this.rightWidth === 100) {
          this.rightHalf();
        } else {
          this.rightFull();
        }
      },

      async fetchTreeChildren(item) {
        return new Promise((resolve) => {
          let that = this;
          this.loadDirectories(item.path, (directories) => {
            if ((directories.length === 1) && (directories[0].name === '/')) {
              delete item.children;
            }
            const key = item.path;
            const parentNode = that.$refs.treeView.nodes[key];

            let childNode;
            directories.forEach((child) => {
              childNode = {
                ...parentNode,
                item: child,
                vnode: null,
              };
              that.$refs.treeView.nodes[child.path] = childNode;
            });

            for (let i = 0; i < directories.length; i += 1) {
              item.children.push(directories[i]);
            }
            resolve();
          });
        })
      },

      countOccurrences(arr, val) {
        return arr.reduce((a, v) => (v === val ? a + 1 : a), 0);
      },

      removeExistingLabels(arr) {
        return this.labels.filter((elements) => !arr.includes(elements));
      },

      getLabels() {
        return
        this.loadingLabels = true;
        const that = this;
        invoke('GET', '/labels', [], (data) => {
          that.labels = data;
          that.loadingLabels = false;
        });
      },

      setLabels(labels, id) {
        let url = '/label?id=' + id + '&label=' + labels;
        invoke('GET', url, [], (data) => {
          this.item.labels = labels;
        });
      },


      setLabelsMultiple(id, add, remove, chunk = false) {
        return
        if (!chunk) {
          this.loadingLabelsOperation = true;
        }
        for (let i = 0; i < add.length; i += 1) {
          this.setLabel(add[i], id)
        }
        for (let i = 0; i < remove.length; i += 1) {
          this.unsetLabel(remove[i], id)
        }
        if (!chunk) {
          this.loadingLabelsOperation = false;
          this.labelsConfirmation = false;
          this.getLabels();
        }
      },

      setLabelsMultipleChunk(items, labels) {
        return
        this.loadingLabelsOperation = true;
        for (let i = 0; i < items.length; i += 1) {
          this.setLabelsMultiple(items[i].id, labels, [], true)
        }
        this.loadingLabelsOperation = false;
        this.labelsConfirmationMultiple = false;
        this.getLabels();
      },

      filterLabels(label) {
        if (label === this.filteredLabel) {
          label = ''
        }

        this.label = label;
        console.log("li5");
        this.loadItems();
        this.filteredLabel = label;
      },

      discardAddLabels(multiple) {
        if (multiple) {
          this.labelsDialogMultiple = false;
          this.labelsToAddMultiple = [];
          this.searchLabelsMultiple = '';
        } else {
          this.labelsDialog = false
          this.labelsToAdd = [];
          this.labelsToRemove = [];
          this.labelsToRemoveElements = [];
          this.searchLabels = '';
        }
      },

      removeFalseLabels(labels) {
        return Object.keys(labels).reduce(function (filtered, key) {
          if (labels[key] === true) {
            filtered[key] = labels[key];
          }
          return filtered;
        }, {});
      },

      loadItems() {
        this.loadElementsTable = true;
        let url = `/items?type=${this.type}`;
        if (!Vue._.isEmpty(this.filter) && this.filter.type !== '' && !Vue._.isEmpty(this.filter.columns)) {
          Vue._.forEach(this.filter.columns, (value, column) => {
            for (let i = 0; i < value.length; i += 1) {
              url += `&filter[${column}]=${encodeURI(value[i])}`;
            }
          });
        }
        if (!Vue._.isEmpty(this.sort) && this.sort.type !== '' && !Vue._.isEmpty(this.sort.columns)) {
          Vue._.forEach(this.sort.columns, (value, column) => {
            let direction = '';
            switch (value) {
              case 'ASC':
                direction = 'ASC';
                break;
              case 'DESC':
              default:
                direction = 'DESC';
            }
            if (direction !== '') {
              url += `&sort[${column}]=${direction}`;
            }
          });
        }
        if (this.label !== '') {
          url += `&labels=${this.label}`;
        }

        url += `&offset=${this.offset}`;
        url += `&limit=${this.limit}`;

        let that = this;
        invoke('GET', url, [], (items) => {
          that.type = this.type;
          that.setItems(items.elements);
          this.loadElementsTable = false;
          that.itemCount = items.count;
        });
      },

      loadDirectories(path, callback) {
        let slash = '';
        let url = '';

        switch (this.type) {
          case 'directory':
            slash = '/';
            break;
          case 'file':
            slash = '/';
            break;
          case 'windows-registry-key':
            slash = '\\';
            break;
          default:
            callback([]);
            return
        }

        if ((this.type === 'windows-registry-key') && (path === '')) {
          url = `/tree?directory=${path}`;
        } else {
          url = `/tree?directory=${path}${slash}`;
        }

        url += `&type=${this.type}`;

        let that = this;
        invoke('GET', url, [], (response) => {
          let directories = [];
          for (let i = 0; i < response.length; i += 1) {
            if (response[i].name === '/') {
              continue
            }
            if ((that.type === 'windows-registry-key') && (path === '')) {
              directories.push(
                {
                  path: response[i],
                  name: response[i],
                  children: [],
                },
              );
            } else {
              directories.push(
                {
                  path: path + slash + response[i],
                  name: response[i],
                  children: [],
                },
              );
            }
          }
          callback(directories);
        });
      },

      setItems(data) {
        const calcHeaders = function () {
          const headers = [];
          const keys = Vue._.keys(data[0]);
          for (let i = 0; i < keys.length; i += 1) {
            if (keys[i] !== "id" && keys[i] !== "type") {
              headers.push({
                text: Vue._.startCase(keys[i]),
                value: keys[i],
              });
            }
          }
          return headers;
        };

        if (data.length !== 0 && 'type' in data[0]) {
          // this.type = data[0].type;
          if (data[0].type in templates && 'headers' in templates[data[0].type]) {
            this.headers = templates[this.type].headers;
          } else {
            this.headers = calcHeaders();
          }
        } else {
          this.headers = calcHeaders();
        }
        this.offset = 0;
        this.items = data;
      },
    },

    mounted() {
      this.loadingType = true;
      this.getDirectories();
      // this.getLabels();
      invoke('GET', '/tables', [], (tables) => {
        this.tables = tables;
        this.loadingType = false;
      });
    },

    destroyed() {
      this.item = {};
    },
  };
</script>

<style scoped>
  table tr td {
    cursor: pointer !important;
  }

  .v-expansion-panel--active > .v-expansion-panel-header {
    min-height: 48px;
  }
</style>

<style lang="scss">
  @import '~vuetify/src/styles/styles.sass';
  @import '../styles/colors.scss';
  @import '../styles/animations.scss';
  @import '~animate.css';
  @import "../styles/items.sass";
</style>
