<template>
  <div
    ref="optionContainer"
    class="content comment-container"
    style="
      border: 1px solid #d0d7de;
      height: 100px;
      border-radius: 6px;
      display: flex;
      flex-direction: column;
      justify-content: space-around;
    "
    :data-discussionId="discussionId"
    :data-codeId="codeId"
    :data-startLine="startLine"
    :data-endLine="endLine"
  >
    <button @click="renderCreateAiCommentForm" class="option-button">
      <span>AI 리뷰 생성</span>
    </button>

    <button @click="renderCreateCommentForm" class="option-button">
      <span>코멘트 작성</span>
    </button>
  </div>
</template>

<script>
import { getComboMarkdownEditor } from "../features/comp/ComboMarkdownEditor";
import { fetchAiSampleCodes } from "../features/repo-ai-samplecode";
import {
  fetchCommentForm,
  initDiscussionFileCommentForm,
} from "./dIscussion-file-comment-form";

export default {
  props: {
    repoLink: {
      type: String,
      required: true,
    },
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

  mounted() {
    setTimeout(() => {
        document.addEventListener("click", this.removeOptionContainer, {
          once: true,
        });
      }, 0);
      
  },

  methods: {
    async renderCreateCommentForm(event) {
      const targetLine = event.target.closest("tr");

      console.log(this.repoLink);
      console.log(this.discussionId);

      const queryParams = {
        discussionId: this.discussionId,
        codeId: this.codeId,
        startLine: this.startLine,
        endLine: this.endLine,
      };

      const requestURL = new URL(`${this.repoLink}/discussions/comment`);
      Object.entries(queryParams).forEach(([key, value]) => {
        requestURL.searchParams.set(key, value);
      });

      const commentForm = await fetchCommentForm(requestURL);

      initDiscussionFileCommentForm(commentForm);

      targetLine.insertAdjacentElement("afterend", commentForm);
      targetLine.remove();
    },

    async renderCreateAiCommentForm(event) {
      const targetLine = event.target.closest("tr");

      const queryParams = {
        discussionId: this.discussionId,
        codeId: this.codeId,
        startLine: this.startLine,
        endLine: this.endLine,
      };

      const requestURL = new URL('/ai/discussion/form', window.location.origin);
      Object.entries(queryParams).forEach(([key, value]) => {
        requestURL.searchParams.set(key, value);
      });

      const commentForm = await fetchCommentForm(requestURL);

      const modalShowBtn = commentForm.querySelector('.show-ai-code-modal')

      modalShowBtn.addEventListener('click', (event) => {
        event.preventDefault()

        const content = commentForm.querySelector("textarea").value;
        
        const data = {
          "content": content,
          "discussionId": this.discussionId,
          "codeId":  Number(this.codeId),
          "startLine": this.startLine,
          "endLine": this.endLine,
        }
        const aiCodeContainers = document.querySelectorAll('.ai-code-area')
        fetchAiSampleCodes(data, aiCodeContainers);
        event.target.closest('tr').remove();
      });
      
      targetLine.insertAdjacentElement("afterend", commentForm);
      targetLine.remove();
    },

    removeOptionContainer(event) {

      if (event.target.closest(".comment-container")) return;


      
      this.$refs.optionContainer.closest("tr").remove();
    }


  },
};
</script>

<style>
.option-button {
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
