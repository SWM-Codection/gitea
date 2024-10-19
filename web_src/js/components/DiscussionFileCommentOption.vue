<template>
  <tr>
    <td colspan="3">
      <div
      ref = "optionContainer"
        class="content comment-container"
        style="border: 1px solid #d0d7de; border-radius: 6px"
        :data-discussionId="discussionId"
        :data-codeId="codeId"
        :data-startLine="startLine"
        :data-endLine="endLine"
      >
        <button class="option-button">
          <span>AI 리뷰 생성</span>
        </button>

        <button @click="renderCreateCommentForm" class="option-button">
          <span>코멘트 작성</span>
        </button>
      </div>
    </td>
  </tr>
</template>

<script>
import {
  fetchCommentForm,
  initDiscussionFileCommentForm,
} from "./dIscussion-file-comment-form";

export default {
  props: {
    discussionId: {
      type: Number,
      required: true,
    },
    codeId: {
      type: Number,
      required: true,
    },
    startLine: {
      type: Number,
      required: true,
    },
    endLine: {
      type: Number,
      required: true,
    },
  },
  methods: {
    async renderCreateCommentForm(event) {
        this.$refs.optionContainer.remove();
        
      // if (!this.isDraggingForComment) {

      //   const { codeId, lineNumber } = this.extractDataFromLine(targetLine);

      //   const codeLinePosition = this.createCodePosition(codeId, lineNumber);
      //   this.currentDraggedRange = this.createCodeLineRange(
      //     codeId,
      //     codeLinePosition,
      //     codeLinePosition,
      //   );
      // }

      const requestURL = new URL(`${this.repoLink}/discussions/comment`);
      Object.entries(queryParams).forEach(([key, value]) => {
        requestURL.searchParams.set(key, value);
      });

      const commentForm = await fetchCommentForm(requestURL);

      initDiscussionFileCommentForm(commentForm);

      targetLine.insertAdjacentElement("afterend", commentForm);
    },
  },
};
</script>

<style>
.option-button {
  position: absolute;
  cursor: pointer;
  color: black;
  background-color: #f1f3f5;
  top: 16px;
  right: 4px;
  padding: 8px;
  z-index: 1;
  border-radius: 6px;
  border: 1px solid grey;
  justify-content: center;
  font-size: smaller;
}
</style>
