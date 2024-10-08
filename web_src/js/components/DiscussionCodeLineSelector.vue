<template>
  <div
    class="file-header ui top attached header tw-items-center tw-justify-between tw-flex-wrap"
    style="position: sticky; top: 0; z-index: 999"
  >
    <div class="file-info tw-font-mono">
      <div :id="`discussion-${content.NameHash}`" class="file-info-entry">
        {{ content.Name }}
      </div>
    </div>
  </div>
  <div class="ui bottom attached table unstackable segment">
    <div class="file-view code-view" style="display: flex">
      <table :id="content.Name" ref="codeTable" class="discussion-file-table">
        <tbody v-for="codeBlock in content.codeBlocks" :key="codeBlock.codeId">
          <DiscussionFileCodeLine
<<<<<<< HEAD
          :lines="codeBlock.lines"
          :codeId="codeBlock.codeId"
          @show-comment-form="showCommentForm"
          @handle-mouse-down="handleMouseDown"
=======
            :lines="codeBlock.lines"
            :codeId="codeBlock.codeId"
            @show-comment-form="renderCreateCommentForm"
            @handle-mouse-down="handleMouseDown"
>>>>>>> 75358a09f8 (main 최신화 (#113))
          />
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
<<<<<<< HEAD

import { GET, POST } from "../modules/fetch";
import {
  initComboMarkdownEditor,
  validateTextareaNonEmpty,
} from "../features/comp/ComboMarkdownEditor";
import DiscussionFileCodeLine from "./DiscussionFileCodeLine.vue";
=======
import { GET, POST } from "../modules/fetch";
import DiscussionFileCodeLine from "./DiscussionFileCodeLine.vue";
import {initDiscussionCommentsEventHandler} from "../features/discussion-file-comment.js";
import { convertTextToHTML, createCommentPlaceHolder, fetchCommentForm, initDiscussionFileCommentForm} from "./dIscussion-file-comment-form.js";
import { initAiSampleCodeModal } from "../features/repo-ai-samplecode.js";
>>>>>>> 75358a09f8 (main 최신화 (#113))

const { pageData } = window.config;

export default {
<<<<<<< HEAD
  components: {  DiscussionFileCodeLine },
=======
  components: { DiscussionFileCodeLine },
>>>>>>> 75358a09f8 (main 최신화 (#113))

  props: {
    content: {
      type: Object,
      required: true,
    },
  },

  data() {
    return {
      currentDraggedPosition: null,
      isDraggingForComment: false,
      currentDraggedRange: null,
      showMultiLineCommentForm: null,
      repoLink: pageData.RepoLink,
      discussionId: pageData.DiscussionId,
      errorText: "",
    };
  },

  async mounted() {
    await this.fetchDiscussionComments();
<<<<<<< HEAD
=======
    initAiSampleCodeModal() 
>>>>>>> 75358a09f8 (main 최신화 (#113))
  },

  methods: {
    isFileSelecting(fileElement) {
      const table = fileElement.closest(".discussion-file-table");
      return table ? table.classList.contains("is-selecting") : false;
    },

    setSelection(target, canExpand) {
<<<<<<< HEAD
      const {codeId, lineNumber} = this.extractDataFromLine(target)

      if (canExpand && this.currentDraggedRange && this.currentDraggedRange.codeId === codeId) {

=======
      const { codeId, lineNumber } = this.extractDataFromLine(target);

      if (
        canExpand &&
        this.currentDraggedRange &&
        this.currentDraggedRange.codeId === codeId
      ) {
>>>>>>> 75358a09f8 (main 최신화 (#113))
        if (lineNumber < this.currentDraggedRange.startPosition.lineNumber) {
          return;
        }

        const expandedRange = this.createCodeLineRange(
          codeId,
          this.currentDraggedRange.startPosition,
<<<<<<< HEAD
          this.createCodePosition(codeId, lineNumber)
        );

        this.showMultiLineCommentForm = () => {
          const button = target.closest("tr")?.querySelector(".add-code-comment");
=======
          this.createCodePosition(codeId, lineNumber),
        );

        this.showMultiLineCommentForm = () => {
          const button = target
            .closest("tr")
            ?.querySelector(".add-code-comment");
>>>>>>> 75358a09f8 (main 최신화 (#113))
          if (button) {
            button.click();
          }
        };
        this.displayHighlight(expandedRange);
      } else {
        const position = this.createCodePosition(codeId, lineNumber);
<<<<<<< HEAD
        const expandedRange = this.createCodeLineRange(codeId, position, position);
=======
        const expandedRange = this.createCodeLineRange(
          codeId,
          position,
          position,
        );
>>>>>>> 75358a09f8 (main 최신화 (#113))
        this.displayHighlight(expandedRange);
      }
    },

    createCodeLineRange(codeId, startPosition, endPosition) {
      return {
        codeId,
        startPosition,
        endPosition,
        elements: () => {
          const startLine = Number(startPosition.lineNumber);
          const endLine = Number(endPosition.lineNumber);
          const lineElements = new Set();

          for (let i = startLine; i <= endLine; i++) {
<<<<<<< HEAD
            const lineElement = this.$refs.codeTable.querySelector(`#line-${codeId}-${i}`);
=======
            const lineElement = this.$refs.codeTable.querySelector(
              `#line-${codeId}-${i}`,
            );
>>>>>>> 75358a09f8 (main 최신화 (#113))
            if (lineElement) {
              lineElements.add(lineElement);
            }
          }

          return lineElements;
        },
      };
    },

    createCodePosition(codeId, lineNumber) {
      return { codeId, lineNumber: Number(lineNumber) };
    },

    displayHighlight(range) {
      if (this.currentDraggedRange) {
        for (const el of this.currentDraggedRange.elements()) {
          el.classList.remove("selected-line");
        }
      }

      this.currentDraggedRange = range;

      for (const el of this.currentDraggedRange.elements()) {
        el.classList.add("selected-line");
      }
    },

    removeHighlight() {
      if (this.currentDraggedRange) {
        for (const el of this.currentDraggedRange.elements()) {
          el.classList.remove("selected-line");
        }
        this.currentDraggedRange = null;
      }
    },

    handleMouseDown(event) {
      if (!(event instanceof MouseEvent) || event.button !== 0) {
        return;
      }

      const targetElement = event.currentTarget;
<<<<<<< HEAD
      const lineNumberElement = this.prevLinkableLine(targetElement.parentElement);
=======
      const lineNumberElement = this.prevLinkableLine(
        targetElement.parentElement,
      );
>>>>>>> 75358a09f8 (main 최신화 (#113))

      if (!lineNumberElement) {
        return;
      }

      const table = this.$refs.codeTable;
      if (!table) {
        return;
      }

<<<<<<< HEAD
      this.addCommentDragSelectionEvent(table);
=======
      this.addCodeDragSelectionEvent(table);
>>>>>>> 75358a09f8 (main 최신화 (#113))
      this.currentDraggedPosition = lineNumberElement;
      this.isDraggingForComment = true;

      const mouseUpHandler = () => {
<<<<<<< HEAD
        this.removeCommentDragSelectionEvent(table);
=======
        this.removeCodeDragSelectionEvent(table);
>>>>>>> 75358a09f8 (main 최신화 (#113))
        this.currentDraggedPosition = null;
        this.isDraggingForComment = false;
      };

      targetElement.addEventListener("mouseup", mouseUpHandler, { once: true });

<<<<<<< HEAD
      if (this.currentDraggedRange && this.currentDraggedRange.elements().size > 1) {
=======
      if (
        this.currentDraggedRange &&
        this.currentDraggedRange.elements().size > 1
      ) {
>>>>>>> 75358a09f8 (main 최신화 (#113))
        event.preventDefault();
      }
    },

<<<<<<< HEAD
    commentDragSelectionIfMouseEnterToCode(codeElement) {
=======
    codeDragSelectionIfMouseEnterToCode(codeElement) {
>>>>>>> 75358a09f8 (main 최신화 (#113))
      const target = this.prevLinkableLine(codeElement);
      if (!target || !this.isFileSelecting(codeElement)) {
        return;
      }
      this.setSelection(target, true);
    },

<<<<<<< HEAD
    commentDragSelectionIfMouseEnterToLineNumber(lineNumberElement) {
      this.setSelection(lineNumberElement, true);
    },

    addCommentDragSelectionEvent(table) {
      table.addEventListener("mouseenter", this.handleDragMouseEvent, { capture: true });
    },

    removeCommentDragSelectionEvent(table) {
      this.isDraggingForComment = false;
      table.removeEventListener("mouseenter", this.handleDragMouseEvent, { capture: true });
      setTimeout(() => {
        document.addEventListener("click", this.handleClickOutside, { once: true });
=======
    codeDragSelectionIfMouseEnterToLineNumber(lineNumberElement) {
      this.setSelection(lineNumberElement, true);
    },

    addCodeDragSelectionEvent(table) {
      table.addEventListener("mouseenter", this.handleDragMouseEvent, {
        capture: true,
      });
    },

    removeCodeDragSelectionEvent(table) {
      this.isDraggingForComment = false;
      table.removeEventListener("mouseenter", this.handleDragMouseEvent, {
        capture: true,
      });
      setTimeout(() => {
        document.addEventListener("click", this.handleClickOutside, {
          once: true,
        });
>>>>>>> 75358a09f8 (main 최신화 (#113))
      }, 0);
    },

    handleDragMouseEvent(event) {
      const target = event.target.closest("tr");
      if (!(target instanceof Element)) {
        return;
      }

      if (this.currentDraggedPosition) {
        this.beginDrag();
      }

      const linesNum = target.querySelector(".lines-num");
      const linesCode = target.querySelector(".lines-code");

      if (linesNum && linesNum.classList.contains("lines-num")) {
<<<<<<< HEAD
        this.commentDragSelectionIfMouseEnterToLineNumber(linesNum);
      } else if (linesCode && linesCode.classList.contains("lines-code")) {
        this.commentDragSelectionIfMouseEnterToCode(linesCode);
=======
        this.codeDragSelectionIfMouseEnterToLineNumber(linesNum);
      } else if (linesCode && linesCode.classList.contains("lines-code")) {
        this.codeDragSelectionIfMouseEnterToCode(linesCode);
>>>>>>> 75358a09f8 (main 최신화 (#113))
      }
    },

    beginDrag() {
      if (!this.currentDraggedPosition) {
        return;
      }

      this.setSelection(this.currentDraggedPosition, false);
      const table = this.$refs.codeTable;
      table.classList.add("is-selecting");
      this.currentDraggedPosition = null;

      const mouseUpHandler = (event) => {
        table.classList.remove("is-selecting", "is-commenting");
        if (this.showMultiLineCommentForm) {
          this.showMultiLineCommentForm();
          this.showMultiLineCommentForm = null;
        }
<<<<<<< HEAD
        this.removeCommentDragSelectionEvent(table);
=======
        this.removeCodeDragSelectionEvent(table);
>>>>>>> 75358a09f8 (main 최신화 (#113))
        event.preventDefault();
      };

      document.addEventListener("mouseup", mouseUpHandler, { once: true });
    },

    prevLinkableLine(element) {
      if (element.classList.contains("lines-num")) {
        return element;
      }
      const prevSibling = element.previousElementSibling;
      return prevSibling ? this.prevLinkableLine(prevSibling) : null;
    },

    handleClickOutside(event) {
      if (!this.currentDraggedRange) {
        return;
      }
      const target = event.target;
      if (target.closest(".discussion-file-table")) {
        return;
      }
      this.removeHighlight();
    },

    extractDataFromLine(line) {
      const targetProperties = line.id.split("-");
      const codeId = targetProperties[1];
      const lineNumber = targetProperties[2];
<<<<<<< HEAD
      return {codeId, lineNumber}
    },

    
    createCommentPlaceHolder(commentText) {
      const placeholder = document.createElement("tr");
      const td = document.createElement("td");
        td.innerHTML = commentText;
        td.setAttribute("colspan", "3");
        placeholder.appendChild(td);
      return placeholder
    },

    async showCommentForm(event) {
      if (!this.isDraggingForComment) {
        const line = event.target.closest("tr");
        
        const {codeId, lineNumber} = this.extractDataFromLine(line)
        
        const codeLinePosition = this.createCodePosition(codeId, lineNumber);
        this.currentDraggedRange = this.createCodeLineRange(codeId, codeLinePosition, codeLinePosition);
=======
      return { codeId, lineNumber };
    },

    async renderCreateCommentForm(event) {
      const targetLine = event.target.closest("tr");
      if (!this.isDraggingForComment) {
        
        const { codeId, lineNumber } = this.extractDataFromLine(targetLine);

        const codeLinePosition = this.createCodePosition(codeId, lineNumber);
        this.currentDraggedRange = this.createCodeLineRange(
          codeId,
          codeLinePosition,
          codeLinePosition,
        );
>>>>>>> 75358a09f8 (main 최신화 (#113))
      }

      const { codeId, startPosition, endPosition } = this.currentDraggedRange;
      const queryParams = {
        discussionId: this.discussionId,
<<<<<<< HEAD
        codeId,
=======
        codeId : codeId,
>>>>>>> 75358a09f8 (main 최신화 (#113))
        startLine: startPosition.lineNumber,
        endLine: endPosition.lineNumber,
      };

      const requestURL = new URL(`${this.repoLink}/discussions/comment`);
      Object.entries(queryParams).forEach(([key, value]) => {
        requestURL.searchParams.set(key, value);
      });

<<<<<<< HEAD
      try {
        const response = await GET(requestURL.toString());
        if (!response.ok) {
          this.errorText = response.statusText;
          return;
        }
        const body = await response.text();

        const placeholder = this.createCommentPlaceHolder(body)

        const targetLine = event.target.closest("tr");
        targetLine.insertAdjacentElement("afterend", placeholder);

        await initComboMarkdownEditor(td.querySelector(".combo-markdown-editor"));
        await this.initDiscussionFileCommentForm(td.querySelector("form"));

        placeholder.addEventListener("click", this.removeCommentForm, { capture: true });
      } catch (err) {
        this.errorText = err.message;
        console.error(this.errorText);
      }
    },

    removeCommentForm(event) {
      if (event.target && event.target.classList.contains("cancel-code-comment")) {
        const commentForm = event.target.closest("tr");
        if (commentForm) {
          commentForm.remove();
        }
      }
    },

    async submitDiscussionFileCommentForm(event) {
      event.preventDefault();
      const form = event.target;

      const textarea = form.querySelector("textarea");
      if (!validateTextareaNonEmpty(textarea)) {
        return;
      }

      if (form.classList.contains("is-loading")) {
        return;
      }

      try {
        form.classList.add("is-loading");
        const formData = new FormData(form);

        const response = await POST(
          `${this.repoLink}/discussions/${this.discussionId}/comment`,
          { data: formData }
        );

        if (!response.ok) {
          this.errorText = response.statusText;
          return;
        }

        const body = await response.json();

        const resp = await GET(`${this.repoLink}/discussions/comment/${body.id}`);
        const commentHolderHTML = await resp.text();

        const tempDiv = document.createElement("div");
        tempDiv.innerHTML = commentHolderHTML;
        const commentHolder = tempDiv.firstElementChild;

        form.closest(".discussion-file-comment-holder").replaceWith(commentHolder);
      } catch (e) {
        this.errorText = e.message;
        console.error(this.errorText);
      } finally {
        form.classList.remove("is-loading");
      }
    },

    convertTextToHTML(text) {
      const tempDiv = document.createElement("div");
      tempDiv.innerHTML = text
      return tempDiv.firstElementChild
=======
      const commentForm = await fetchCommentForm(requestURL)  

      initDiscussionFileCommentForm(commentForm);

      targetLine.insertAdjacentElement("afterend", commentForm);
>>>>>>> 75358a09f8 (main 최신화 (#113))
    },

    async fetchDiscussionComments() {
      try {
        const codeBlocks = this.content.codeBlocks;
<<<<<<< HEAD
        const commentPromises = codeBlocks.flatMap((codeBlock) => {
          const { codeId, comments } = codeBlock;
          return comments.map(async (comment) => {
            const response = await GET(`${this.repoLink}/discussions/comment/${comment.id}`);
            const result = await response.text();

            const commentHolder = this.convertTextToHTML(result)

            return { comment, commentHolder, codeId };
          });
        });

        const allComments = await Promise.all(commentPromises);

        allComments.forEach(({ comment, commentHolder, codeId }) => {
          const targetLine = this.$refs.codeTable.querySelector(`#line-${codeId}-${comment.endLine}`);
=======
        const commentPromises = codeBlocks.map(async (codeBlock) => {
          const { codeId } = codeBlock;
          const response = await GET(
            `${this.repoLink}/discussions/comments/${codeId}`,
          );

          const commentGroups = await response.json();

          return commentGroups.map((result) => {
            const line = result.endLine;
            const commentHolder = result.html;

            return { line, commentHolder, codeId };
          });
        });

        const allCommentsNested = await Promise.all(commentPromises);
        const allComments = allCommentsNested.flat();

        allComments.forEach(({ line, commentHolder, codeId }) => {
          const targetLine = this.$refs.codeTable.querySelector(
            `#line-${codeId}-${line}`,
          );
          commentHolder = convertTextToHTML(commentHolder)
>>>>>>> 75358a09f8 (main 최신화 (#113))

          if (targetLine) {
            const tr = document.createElement("tr");
            const td = document.createElement("td");
            td.setAttribute("colspan", "3");
            td.appendChild(commentHolder);
<<<<<<< HEAD
=======
            initDiscussionCommentsEventHandler(commentHolder);
>>>>>>> 75358a09f8 (main 최신화 (#113))
            tr.appendChild(td);
            targetLine.insertAdjacentElement("afterend", tr);
          }
        });
      } catch (e) {
        console.error("Error processing code blocks:", e);
      }
    },

<<<<<<< HEAD
    async initDiscussionFileCommentForm(form) {
      form.addEventListener("submit", this.submitDiscussionFileCommentForm);
    },
=======

>>>>>>> 75358a09f8 (main 최신화 (#113))
  },
};
</script>

<style scoped>
.selected-line {
  background-color: #f5f5dc;
}
.is-selecting {
  cursor: pointer;
}
</style>
