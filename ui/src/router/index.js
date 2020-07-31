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
import VueRouter from 'vue-router';
import items from '../views/Items.vue';
import file from '../views/File.vue';
import files from '../views/Files.vue';
import tasks from '../views/Tasks.vue';
import logs from '../views/Logs.vue';
import errors from '../views/Errors.vue';

Vue.use(VueRouter);

const routes = [
  {
    path: '/',
    redirect: '/items',
  },
  {
    path: '/tasks',
    name: 'tasks',
    component: tasks,
  },
  {
    path: '/items',
    name: 'items',
    component: items,
  },
  {
    path: '/files',
    name: 'files',
    component: files,
  },
  {
    path: '/logs',
    name: 'logs',
    component: logs,
  },
  {
    path: '/errors',
    name: 'errors',
    component: errors,
  },
  {
    path: '/file',
    name: 'file',
    component: file,
  },
  {
    path: '/workflows',
    name: 'workflows',
    component: () => import('../views/Workflows.vue'),
  },
];

const router = new VueRouter({
  routes,
});

export default router;
