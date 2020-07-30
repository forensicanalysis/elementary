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
  <div class="task scrollableArea pa-8 elevation-3"
       style="height: 100%; width: 1024px; margin: 40px auto; background-color: #eee;">
    <v-row dense class="" v-for="(task, id) in tasks" :key="id"
           :class="{ inactive: (!active[task.Name] || _.isEmpty(cleanSchema(task.Schema).properties)) }">
      <v-col cols="6" class="py-0">
        <v-checkbox
          v-model="active[task.Name]"
          :label="task.Description"
        />
      </v-col>
      <v-col cols="6" v-if="active[task.Name]" class="py-0"  @click="click(task.Name, $event)">
        <v-form ref="form" v-model="valids[task.Name]">
          <v-jsf v-model="models[task.Name]" :schema="cleanSchema(task.Schema)"
                 :options="options" />
        </v-form>
      </v-col>
    </v-row>

    <v-row class="justify-center pa-5">
      <v-btn x-large color="primary" @click="run">
        <v-icon>mdi-play</v-icon>
        Run
      </v-btn>
    </v-row>
    <v-overlay
      :absolute="true"
      :opacity="0.5"
      :value="overlay"
      style="text-align: center"
    >
      <p>Running tasks</p>
      <v-progress-circular indeterminate size="64"></v-progress-circular>
    </v-overlay>

  </div>
</template>

<script>
  import {invoke} from "../store/invoke";
  import VJsf from '@koumoul/vjsf';


  export default {
    name: 'tasks',
    components: {
      VJsf
    },
    data() {
      return {
        overlay: false,
        tasks: [],
        active: {},
        valids: {},
        models: {},
        options: {
          debug: false,
          disableAll: false,
          autoFoldObjects: true
        },
      }
    },

    methods: {
      loadTasks() {
        invoke('GET', '/tasks', [], (tasks) => {
          this.tasks = tasks;
        });
      },

      cleanSchema(schema) {
        delete schema["properties"]["add-to-store"];
        delete schema["properties"]["format"];
        delete schema["properties"]["output"];
        delete schema["properties"]["filter"]; // TODO
        return schema
      },

      click(task, e) {
        if (e.path[0].tagName === "INPUT" && e.path[0].id === "file") {
          e.preventDefault();

          if ('electron' in window) { // electron check
            window.electron.showOpenDialog().then(result => {
              if (!result.canceled) {
                let id = e.path[0].id
                this.$set(this.models[task], id, result.filePaths[0]);
              }
            }).catch(err => {
              console.log(err)
            })
          }

          // TODO: handle files in server mode
        }
      },

      run() {
        this.overlay = true;
        let that = this;
        that._.forEach(this.tasks, function (task, id) {
            if (that.active[task.Name]) {
              invoke('POST', '/run?name=' + task.Name, that.models[task.Name], (result) => {
                that.overlay = false;
              });
            }
          }
        );
      },
    },

    mounted() {
      this.loadTasks();
    },
  };
</script>

<style>
  .task > .row {
    height: 70px;
    transition: height 0.1s ease-in;
  }

  .task > .row.inactive {
    /* border-top: 1px solid rgba(0, 0, 0, 0.12); */
    height: 28px;
  }

  .task .v-input .v-label {
    font-size: 18px;
  }

  .container {
    background-color: #eeeeee;
    margin: 40px 140px;
  }

  .task .vjsf-property .row {
    flex-direction: row;
    flex-wrap: nowrap;
  }

  .task .col-2 label {
    /* font-weight: bold; */
    color: #333;
    font-size: 0.8rem;
  }

  .task .v-input--checkbox {
    margin-top: 0;
  }
</style>
