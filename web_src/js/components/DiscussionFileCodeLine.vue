<template>
  <tr
    v-for="line in lines"
    :id="`line-${this.codeId}-${line.lineNumber}`"
    class="code-line"
    :key="`${codeId}-${line.lineNumber}`"
  >
    <td class="lines-num" :id="`num-${this.codeId}-${line.lineNumber}`">
      {{ line.lineNumber }}
    </td>
    <td>
      <button
        @mousedown="onMouseDown($event)"
        @click="$emit('show-comment-form', $event)"
        class="ui primary button add-code-comment add-code-comment-right"
      >
        <SvgIcon name="octicon-plus" />
      </button>
    </td>
    <td class="lines-code chroma">
      {{ line.content }}
    </td>
  </tr>
</template>

<script>
import { SvgIcon } from '../svg';

export default {
  name: 'CodeLine',
  components : {SvgIcon},
  props: {
    lines: {
      type: Array,
      required: true,
    },
    codeId: {
      type: String,
      required: true,
    },
  },
  computed: {
  },
  methods: {
    onMouseDown(event) {
      // `line` 데이터를 함께 전달합니다.
      this.$emit('handle-mouse-down', event);
    },
  },
};
</script>

<style scoped>

.selected-line {
  background-color: #f5f5dc;
}
</style>