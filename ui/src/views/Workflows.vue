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
  <div class="content">
    <svg
      class="w-full p-5"
      :width="width"
      :height="height"
      :viewBox="viewBox"
      :id="id"
    />
    <SimpleFlowchart :scene.sync="data"/>
  </div>
</template>

<script>
import * as d3Base from 'd3';
import * as d3Dag from 'd3-dag';
import SimpleFlowchart from 'vue-simple-flowchart';

import 'vue-simple-flowchart/dist/vue-flowchart.css';

const d3 = Object.assign({}, d3Base, d3Dag);

export default {
  name: 'workflows',
  components: {
    SimpleFlowchart,
  },
  data() {
    return {
      data: {
        centerX: 1024,
        centerY: 140,
        scale: 1,
        nodes: [
          {
            id: 2,
            x: -700,
            y: -69,
            type: 'Action',
            label: 'test1',
          },
          {
            id: 4,
            x: -357,
            y: 80,
            type: 'Script',
            label: 'test2',
          },
          {
            id: 6,
            x: -557,
            y: 80,
            type: 'Rule',
            label: 'test3',
          },
        ],
        links: [
          {
            id: 3,
            from: 2, // node id the link start
            to: 4, // node id the link end
          },
        ],
      },
      boxW: 150,
      boxH: 60,
      nodeRadius: 80,
      width: 1000,
      height: 250,
      id: 'asd',
      canvas: null,
    };
  },
  computed: {
    viewBox() {
      let vb = `${-this.nodeRadius} ${-this.nodeRadius}`;
      vb += ` ${this.width + 2 * this.nodeRadius}`;
      vb += ` ${this.height + 2 * this.nodeRadius}`;
      return vb;
    },
    graph() {
      const flat = [];
      const { tasks } = this.$store.state;
      for (const name in tasks) {
        if (tasks.hasOwnProperty(name)) {
          const task = tasks[name];
          if ('requires' in task) {
            task.requires.forEach((requirement) => {
              flat.push([requirement, name]);
            });
          }
        }
      }
      return flat;
    },
  },
  mounted() {
    const canvas = d3.select(`#${this.id}`);

    this.daggy(canvas, this.graph, this.nodeRadius, this.width, this.height);

    let elem = d3.select('#fetch');
    elem.attr('fill', '#27ae60');
    elem = d3.select('#unpack');
    elem.attr('fill', '#27ae60');
    elem = d3.select('#plaso');
    elem.attr('fill', '#f39c12');
    elem = d3.select('#hash');
    elem.attr('fill', '#c0392b');
  },
  methods: {
    daggy(svgSelection, idata, nodeRadius, height, width) {
      const defs = svgSelection.append('defs'); // For gradients

      const dag = d3.dagConnect()(idata);

      // let graph = d3.zherebko().size([width, height])(dag);
      // const graph = d3
      //   .sugiyama()
      //   .size([width, height])
      //   .layering(d3.layeringCoffmanGraham())
      //   .decross(d3.decrossOpt())
      //   .coord(d3.coordVert())(dag);

      // const steps = dag.size();
      // const interp = d3.interpolateGreys; // interpolateRainbow
      const colorMap = {};
      dag.each((node, i) => {
        colorMap[node.id] = '#7f8c8d'; // interp(i / steps);
      });

      // How to draw edges
      const line = d3
        .line()
        .curve(d3.curveMonotoneX)
        .x(d => d.y)
        .y(d => d.x);

      // Plot edges
      svgSelection
        .append('g')
        .selectAll('path')
        .data(dag.links())
        .enter()
        .append('path')
        .attr('d', ({ data }) => line(data.points))
        .attr('fill', 'none')
        .attr('stroke-width', 3)
        .attr('stroke', ({ source, target }) => {
          const gradId = `${source.id}-${target.id}`;
          const grad = defs
            .append('linearGradient')
            .attr('id', gradId)
            .attr('gradientUnits', 'userSpaceOnUse')
            .attr('x1', source.y)
            .attr('x2', target.y)
            .attr('y1', source.x)
            .attr('y2', target.x);
          grad
            .append('stop')
            .attr('offset', '0%')
            .attr('stop-color', colorMap[source.id]);
          grad
            .append('stop')
            .attr('offset', '100%')
            .attr('stop-color', colorMap[target.id]);

          return `url(#${gradId})`;
        });

      // Select nodes
      const nodes = svgSelection
        .append('g')
        .selectAll('g')
        .data(dag.descendants())
        .enter()
        .append('g')
        .attr('transform', ({ x, y }) => `translate(${y}, ${x})`);

      // Plot node circles
      nodes
        .append('rect')
        .attr('id', n => n.id)
        .attr('x', -(this.boxW / 2))
        .attr('y', -(this.boxH / 2))
        .attr('rx', '10')
        .attr('ry', '10')
        .attr('width', this.boxW)
        .attr('height', this.boxH)
        .attr('fill', n => colorMap[n.id]);

      const arrow = d3
        .symbol()
        .type(d3.symbolTriangle)
        .size((nodeRadius * nodeRadius) / 20.0);
      svgSelection
        .append('g')
        .selectAll('path')
        .data(dag.links())
        .enter()
        .append('path')
        .attr('d', arrow)
        .attr('transform', ({ source, target, data }) => {
          const [end, start] = data.points.reverse();
          // This sets the arrows the node radius (20) + a little bit (3) away from the node center,
          // on the last line segment of the edge. This means that edges that only span ine level
          // will work perfectly, but if the edge bends, this will be a little off.
          const dx = start.y - end.y;
          const dy = start.x - end.x;
          const scale = (nodeRadius * 1.25) / Math.sqrt(dx * dx + dy * dy);
          // This is the angle of the last line segment
          const angle = ((Math.atan2(-dy, -dx) * 180 * 1.2) / Math.PI + 90) * 1;
          return `translate(${end.y + dx * scale + 10}, ${end.x
            + dy * scale}) rotate(${angle})`;
        })
        .attr('fill', ({ target }) => colorMap[target.id])
        .attr('stroke', 'white')
        .attr('stroke-width', 1.5);

      // Add text to nodes
      nodes
        .append('text')
        .text(d => d.id)
        .attr('font-weight', 'bold')
        .attr('font-family', 'sans-serif')
        .attr('text-anchor', 'middle')
        .attr('alignment-baseline', 'middle')
        .attr('fill', 'white')
        .attr('font-size', '20');

      // return svgNode;
    },
  },
};
</script>
