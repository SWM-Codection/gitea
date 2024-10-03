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
            :lines="codeBlock.lines"
            :codeId="codeBlock.codeId"
            @show-comment-form="renderCreateCommentForm"
            @handle-mouse-down="handleMouseDown"
          />
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import { GET, POST } from "../modules/fetch";
import DiscussionFileCodeLine from "./DiscussionFileCodeLine.vue";
import {initDiscussionCommentsEventHandler} from "../features/discussion-file-comment.js";
import { convertTextToHTML, createCommentPlaceHolder, fetchCommentForm, initDiscussionFileCommentForm} from "./dIscussion-file-comment-form.js";
import { initAiSampleCodeModal } from "../features/repo-ai-samplecode.js";

const { pageData } = window.config;

export default {
  components: { DiscussionFileCodeLine },

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
    initAiSampleCodeModal() 
  },

  methods: {
    isFileSelecting(fileElement) {
      const table = fileElement.closest(".discussion-file-table");
      return table ? table.classList.contains("is-selecting") : false;
    },

    setSelection(target, canExpand) {
      const { codeId, lineNumber } = this.extractDataFromLine(target);

      if (
        canExpand &&
        this.currentDraggedRange &&
        this.currentDraggedRange.codeId === codeId
      ) {
        if (lineNumber < this.currentDraggedRange.startPosition.lineNumber) {
          return;
        }

        const expandedRange = this.createCodeLineRange(
          codeId,
          this.currentDraggedRange.startPosition,
          this.createCodePosition(codeId, lineNumber),
        );

        this.showMultiLineCommentForm = () => {
          const button = target
            .closest("tr")
            ?.querySelector(".add-code-comment");
          if (button) {
            button.click();
          }
        };
        this.displayHighlight(expandedRange);
      } else {
        const position = this.createCodePosition(codeId, lineNumber);
        const expandedRange = this.createCodeLineRange(
          codeId,
          position,
          position,
        );
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
            const lineElement = this.$refs.codeTable.querySelector(
              `#line-${codeId}-${i}`,
            );
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
      const lineNumberElement = this.prevLinkableLine(
        targetElement.parentElement,
      );

      if (!lineNumberElement) {
        return;
      }

      const table = this.$refs.codeTable;
      if (!table) {
        return;
      }

      this.addCodeDragSelectionEvent(table);
      this.currentDraggedPosition = lineNumberElement;
      this.isDraggingForComment = true;

      const mouseUpHandler = () => {
        this.removeCodeDragSelectionEvent(table);
        this.currentDraggedPosition = null;
        this.isDraggingForComment = false;
      };

      targetElement.addEventListener("mouseup", mouseUpHandler, { once: true });

      if (
        this.currentDraggedRange &&
        this.currentDraggedRange.elements().size > 1
      ) {
        event.preventDefault();
      }
    },

    codeDragSelectionIfMouseEnterToCode(codeElement) {
      const target = this.prevLinkableLine(codeElement);
      if (!target || !this.isFileSelecting(codeElement)) {
        return;
      }
      this.setSelection(target, true);
    },

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
        this.codeDragSelectionIfMouseEnterToLineNumber(linesNum);
      } else if (linesCode && linesCode.classList.contains("lines-code")) {
        this.codeDragSelectionIfMouseEnterToCode(linesCode);
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
        this.removeCodeDragSelectionEvent(table);
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
      }

      const { codeId, startPosition, endPosition } = this.currentDraggedRange;
      const queryParams = {
        discussionId: this.discussionId,
        codeId : codeId,
        startLine: startPosition.lineNumber,
        endLine: endPosition.lineNumber,
      };

      const requestURL = new URL(`${this.repoLink}/discussions/comment`);
      Object.entries(queryParams).forEach(([key, value]) => {
        requestURL.searchParams.set(key, value);
      });

      const commentForm = await fetchCommentForm(requestURL)  

      initDiscussionFileCommentForm(commentForm);

      targetLine.insertAdjacentElement("afterend", commentForm);
    },

    async fetchDiscussionComments() {
      try {
        const codeBlocks = this.content.codeBlocks;
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

          if (targetLine) {
            const tr = document.createElement("tr");
            const td = document.createElement("td");
            td.setAttribute("colspan", "3");
            td.appendChild(commentHolder);
            initDiscussionCommentsEventHandler(commentHolder);
            tr.appendChild(td);
            targetLine.insertAdjacentElement("afterend", tr);
          }
        });
      } catch (e) {
        console.error("Error processing code blocks:", e);
      }
    },


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
