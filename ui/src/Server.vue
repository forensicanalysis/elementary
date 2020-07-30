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

Author(s): Jonas Plum, Kadir Aslan
-->

<template>
  <v-app :style="{background: $vuetify.theme.themes[theme].background}">
    <v-app-bar color="sidebar"
               app
               flat
               dark
               dense
               extended>
      <v-toolbar-title class="mr-8 ml-4">Elementary</v-toolbar-title>
      <v-spacer/>
      <v-toolbar-title class="mr-8 ml-4 body-1">
        md1rejuc_2020-05-11T14-05-57.forensicstore
      </v-toolbar-title>
      <v-spacer/>
      <v-toolbar-items>
        <v-btn class="primary--text" text>Cases</v-btn>
        <v-btn text>Upload</v-btn>
      </v-toolbar-items>
      <v-btn class="mr-3" icon @click.stop="$vuetify.theme.dark = !$vuetify.theme.dark">
        <v-icon>mdi-invert-colors</v-icon>
      </v-btn>
    </v-app-bar>
    <v-content>
      <v-card class="mx-3 mb-3" style="margin-top: -48px; z-index: 10">
        <app/>
      </v-card>
    </v-content>
  </v-app>
</template>

<script>
  import Vuetify from './plugins/vuetify.js';

  export default {
    name: 'app',
    components: {},
    data() {
      return {
      };
    },
    computed: {
      theme() {
        return (Vuetify.framework.theme.dark) ? 'dark' : 'light';
      },
    },
    created() {
    },
  };
</script>

<style>

  body {
    font-family: 'Roboto', sans-serif;
    background-color: #B0BEC5;
  }

  h1, h2, h3, .v-toolbar__title, .title {
    font-family: 'Roboto-Slab', serif;
  }

  .v-toolbar .v-btn__content {
    font-weight: 400 !important;
  }

  .v-btn--fab {
    border-radius: 0;
  }

  .v-toolbar__content {
    /* border-bottom: 1px solid #B0BEC5; */
  }

  .theme--dark.v-list-item--active:hover::before, .theme--dark.v-list-item--active::before,
  .theme--dark.v-treeview .v-treeview-node__root.v-treeview-node--active:hover::before, .theme--dark.v-treeview .v-treeview-node__root.v-treeview-node--active::before {
    opacity: 0;
  }

  dt {
    font-weight: bold;
  }

  dd {
    margin-left: 12px;
  }

  .v-toolbar__content, .v-toolbar__extension {
    padding: 0;
  }

  .primary--text {
    color: #d9045d;
  }

</style>
