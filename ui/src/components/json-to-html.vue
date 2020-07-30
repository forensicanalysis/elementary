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
<script type="text/jsx">
/* eslint-disable */
  export default {
    name: "json-to-html",
    props: {
      json: {
        type: Object,
        required: true,
      },
      bem: {
        type: Boolean,
        default: false,
      },
      rootClassName: {
        type: String,
        default: undefined,
      },
      rootTagName: {
        type: String,
        default: 'dl',
      },
      inheritClassName: {
        type: Boolean,
        default: true,
      },
    },

    created() {
      if (this.bem && !this.inheritClassName) {
        console.warn('[Warn]: JsonToDo, conflict props: For using BEM in prop "inheritClassName" ' +
          'value must be "true", you pass: ' + this.inheritClassName);
      }
    },

    /**
     * This function render DOM three from props data
     * DOM will rendered such as json structure
     * Each elements will have necessary classNames from json keys
     **/
    render(h) {
      const json = this.json;

      const getClassName = ({ parentClassName, className, isElement }) => {
        const inheritClassName = this.inheritClassName;
        if (!inheritClassName) return className;


        const isBem = inheritClassName && this.bem;
        const isNumberedListItem = typeof className === 'number';
        const sep = isBem
          ? isElement ? '_' : '__'
          : '-';
        const prefixClassName = parentClassName
          ? parentClassName.split(' ')[0] + sep // get only first class name
          : '';

        return isNumberedListItem
          ? parentClassName + sep + (className + 1)
          : prefixClassName + className

      };

      const jsToDOM = ({ data, className, parentClassName }) => {
        const isString = typeof data === 'string' || typeof data === 'number';
        const isArray = Array.isArray(data);
        const isObject = typeof data === 'object' && !isArray;
        let newClassName;

        /**
         *  return div single element
         */
        if (isString) {
          newClassName = getClassName({
            parentClassName,
            className,
            isElement: true,
          });

          return <div class={newClassName} ><dt>{newClassName}</dt><dd domPropsInnerHTML={data}/></div>

        }

        /**
         * call jsToDOM() for each element in array/object recursive
         */
        newClassName = getClassName({ parentClassName, className });

        if (isArray) {
          return <dl class={newClassName}><dt>{newClassName}</dt><dd>
            {data.map((item, index) => jsToDOM({
                data: item,
                parentClassName: getClassName({ parentClassName, className }),
                className: index,
              }),
            )}
          </dd></dl>;
        }

        if (isObject) {
          return <dl class={newClassName}><dt>{newClassName}</dt><dd>
            {Object.keys(data).map(key => jsToDOM({
                data: data[key],
                parentClassName: getClassName({ parentClassName, className }),
                className: key,
              }),
            )}
          </dd></dl>
        }
      };


      /**
       * Create root element and call jsToDOM() for each json element
       */
      const root = {
        tag: this.rootTagName,
        class: this.rootClassName,
      };

      const dom = Object.keys(json).map(item => jsToDOM({
        data: json[item],
        className: item,
        parentClassName: root.class,
      }));

      dom.push(this.$slots.default);

      return h(root.tag, { class: root.class }, dom);
    },
  };
</script>
