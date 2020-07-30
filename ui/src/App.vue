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
  <v-app :style="{background: $vuetify.theme.themes['light'].background}">
    <v-tabs class="main-bar" style="border-bottom: 1px solid rgba(0, 0, 0, 0.12); flex: 0;">
      <v-tab v-for="(menuitem, idx) in menu" :key="menuitem.name"
             @click="$router.push({ name: menuitem.route }).catch(err => {})"
             :class="{ 'primary--text' : $router.currentRoute.name === menuitem.route}"
             text>
        <v-icon style="opacity: 0.3" class="mr-2" small v-text="'mdi-'+menuitem.icon"/>
        {{menuitem.name}}
      </v-tab>
      <v-spacer/>
      <v-tab v-for="(menuitem, idx) in rmenu" :key="menuitem.name"
             @click="$router.push({ name: menuitem.route }).catch(err => {})">
        <v-icon style="opacity: 0.3" class="mr-2" small v-text="'mdi-'+menuitem.icon"/>
        {{ menuitem.name }}
      </v-tab>
    </v-tabs>
    <router-view/>
  </v-app>
</template>

<script>
  export default {
    name: 'app',
    data() {
      return {
        menu: [
          {name: 'Elements', route: 'items', icon: 'periodic-table'},
          {name: 'Files', route: 'files', icon: 'file-tree'},
          {name: 'Tasks', route: 'tasks', icon: 'checkbox-marked-outline'},
        ],
        rmenu: [
          {name: 'Errors', route: 'errors', icon: 'alert'},
          {name: 'Logs', route: 'logs', icon: 'script-text-outline'},
        ],
      };
    },
  };
</script>

<style lang="scss">

  @import '~vuetify/src/styles/styles.sass';
  @import './styles/colors.scss';
  @import './styles/animations.scss';
  @import '~animate.css';

  html {
    overflow: hidden;
  }

  body {
    /* font-family: 'Roboto', sans-serif; */
    user-select: none; /* TODO: disable in server mode */
    font-family: caption, sans-serif;
    background-color: #fff; /*#B0BEC5;*/
  }

  .v-application--wrap {
    height: 100vh;
  }

  .scrollableArea {
    overflow: auto !important;
  }

  /*
  ::-webkit-scrollbar {
    width: 5px;
    height: 5px;
  }
  ::-webkit-scrollbar-track {
    background: none;
    border-radius: 10px;
  }
  ::-webkit-scrollbar-track-piece {
    background: none;
    border-radius: 10px;
  }
  ::-webkit-scrollbar-thumb {
    background-color: $c-pink;
    height: 50px;
    border-radius: 0px;
    border: none;
  }
  ::-webkit-scrollbar-corner {
    background: none;
  }
  */

  h1, h2, h3, .v-toolbar__title, .title {
    font-family: 'Roboto-Slab', serif;
  }

  .v-toolbar .v-btn__content {
    font-weight: 400 !important;
  }

  .theme--light.v-tabs.main-bar > .v-tabs-bar {
    background-color: #e8e8e8;
  }

  * {
    border-radius: 0 !important;
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
