<script>
import DiscussionFileAddCommentButton from "./DiscussionFileAddCommentButton.vue";
import { SvgIcon } from "../svg";
import { GET, POST} from "../modules/fetch";
import {
  initComboMarkdownEditor,
  validateTextareaNonEmpty,
} from "../features/comp/ComboMarkdownEditor";
import { discussionResponseDummy } from "../modules/stores";

const { pageData } = window.config;

export default {
  components: { DiscussionFileAddCommentButton, SvgIcon },
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
      await this.fetchDiscussionComment()
    },
  methods: {
    isFileSelecting(fileElement) {
      discussionResponseDummy;
      return fileElement
        .closest(".discussion-file-table")
        ?.classList.contains("is-selecting");
    },
    setSelection(target, canExpand) {
      const targetLineData = target.id.split("-");

      if (canExpand) {
        const codeId = targetLineData[1];
        const endLineNumber = targetLineData[2];

        if (
          this.currentDraggedRange &&
          this.currentDraggedRange.codeId === codeId
        ) {
          if (
            codeId === this.currentDraggedRange.codeId &&
            endLineNumber < this.currentDraggedRange.startPosition.lineNumber
          ) {
            return;
          }

          const expandedRange = this.createCodeLineRange(
            codeId,
            this.currentDraggedRange.startPosition,
            this.createCodePosition(codeId, endLineNumber),
          );

          this.showMultiLineCommentForm = () => {
            const button = target
              .closest("tr")
              ?.querySelector(".add-code-comment");
            if (button && expandedRange) {
              button.click();
            }
          };
          this.displayHighlight(expandedRange);
        }
      } else {
        const codeId = targetLineData[1];
        const endLineNumber = targetLineData[2];
        const expandedRange = this.createCodeLineRange(
          codeId,
          this.createCodePosition(codeId, endLineNumber),
          this.createCodePosition(codeId, endLineNumber),
        );
        this.displayHighlight(expandedRange);
      }
    },

    createCodeLineRange(codeId, startPosition, endPosition) {
      return {
        startPosition: startPosition,
        endPosition: endPosition,
        codeId: codeId,
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
      return {
        codeId: codeId,
        lineNumber: lineNumber,
      };
    },

    displayHighlight(expandedRange) {
      if (this.currentDraggedRange) {
        for (const el of this.currentDraggedRange.elements()) {
          el.classList.remove("selected-line");
        }
        this.currentDraggedRange = null;
      }

      this.currentDraggedRange = this.createCodeLineRange(
        expandedRange.codeId,
        expandedRange.startPosition,
        expandedRange.endPosition,
      );

      for (const el of this.currentDraggedRange.elements()) {
        el.classList.add("selected-line");
      }
    },

    removeHighlight() {
      const dummyRange = this.createCodeLineRange(
        "0",
        this.createCodePosition("0", "0"),
        this.createCodePosition("0", "0"),
      );
      this.displayHighlight(dummyRange);
    },

    handleMouseDown(event) {
      if (!(event instanceof MouseEvent)) {
        return;
      }

      if (event.button !== 0) {
        return;
      }

      const targetElement = event.currentTarget;
      const lineNumber = this.prevLinkableLine(targetElement.parentElement);

      if (!lineNumber) {
        return;
      }

      const table = this.$refs.codeTable;
      if (!table) {
        return;
      }

      this.addCommentDragSelectionEvent(table);
      this.currentDraggedPosition = lineNumber;
      this.isDraggingForComment = true;

      targetElement?.addEventListener(
        "mouseup",
        () => {
          this.removeCommentDragSelectionEvent(table);
          this.currentDraggedPosition = null;
          this.isDraggingForComment = false;
        },
        { once: true },
      );

      if (
        this.currentDraggedRange &&
        this.currentDraggedRange.elements.size > 1
      ) {
        event.preventDefault();
      }
    },

    commentDragSelectionIfMouseEnterToCode(codeElement) {
      const target = prevLinkableLine(codeElement);

      if (!target || !isFileSelecting(codeElement)) {
        return;
      }

      this.setSelection(target, true);
    },

    commentDragSelectionIfMouseEnterToLineNumber(lineNumberElement) {
      this.setSelection(lineNumberElement, true);
    },
    addCommentDragSelectionEvent(table) {
      table.addEventListener("mouseenter", this.handleDragMouseEvent, {
        capture: true,
      });
    },

    removeCommentDragSelectionEvent(table) {
      this.isDraggingForComment = false;
      table.removeEventListener("mouseenter", this.handleDragMouseEvent, {
        capture: true,
      });
      setTimeout(() => {
        document.addEventListener("click", this.handleClick, { once: true });
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

      if (linesNum && linesNum?.classList.contains("lines-num")) {
        this.commentDragSelectionIfMouseEnterToLineNumber(linesNum);
      } else if (linesCode && linesCode?.classList.contains("lines-code")) {
        this.commentDragSelectionIfMouseEnterToCode(linesCode);
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

      document.addEventListener("mouseup", (event) => {
        table.classList.remove("is-selecting", "is-commenting");
        this.showMultiLineCommentForm && this.showMultiLineCommentForm();
        this.showMultiLineCommentForm = null;
        this.removeCommentDragSelectionEvent(table);
        event.preventDefault();
      });
    },

    prevLinkableLine(element) {
      if (element.classList.contains("lines-num")) {
        return element;
      }

      const previousElementSibling = element.previousElementSibling;
      if (previousElementSibling) {
        return this.prevLinkableLine(previousElementSibling);
      }

      return null;
    },

    handleClick(event) {
      if (!this.currentDraggedRange) {
        return;
      }

      const target = event.target;
      if (target?.closest(".discussion-file-table")) {
        return;
      }

      this.removeHighlight();
    },

    async showCommentForm(event) {
      const draggedRange = this.currentDraggedRange;

      const queryParams = {
        discussionId: this.discussionId,
        // TODO: 한 번도 드래그 안한 상태면 빈 값이 들어가서 오류 발생하는 거 해결하기
        codeId: draggedRange.codeId,
        startLine: draggedRange.startPosition.lineNumber,
        endLine: draggedRange.endPosition.lineNumber,
      };

      let requestURL = new URL(`${this.repoLink}/discussions/comment`);

      for (const [key, value] of Object.entries(queryParams)) {
        requestURL.searchParams.set(key, value);
      }

      try {
        let response;

        response = await GET(requestURL.toString());

        if (response.ok) {
          const body = await response.text();
          console.log(body);

          const placeholder = document.createElement("tr");
          const td = document.createElement("td");
          td.innerHTML = body;
          td.setAttribute("colspan", "3");
          placeholder.appendChild(td);

          const targetLine = event.target.closest("tr");
          targetLine.insertAdjacentElement("afterend", placeholder);

          await initComboMarkdownEditor(
            td.querySelector(".combo-markdown-editor"),
          );
          // TODO: 폼 등록
          await this.initDiscussionFileCommentForm(td.querySelector("form"));

          placeholder.addEventListener("click", this.removeCommentForm, {
            capture: true,
          });
        } else {
          this.errorText = response.statusText;
        }
      } catch (err) {
        this.errorText = err.message;
      } finally {
        if (this.errorText) console.log(this.errorText);
        // 에러 메시지 표시 관련 로직 추가
      }
    },

    removeCommentForm: (event) => {
      if (
        event.target &&
        event.target.classList.contains("cancel-code-comment")
      ) {
        const commentForm = event.target.closest("tr");
        if (!commentForm) {
          return;
        }
        commentForm.remove();
      }
    },

    async submitDiscussionFileCommentForm(event) {
      if (!event.target) {
        return;
      }
      event.preventDefault();

      const form = event.target;
      const textarea = form.querySelector("textarea");
      if (!validateTextareaNonEmpty(textarea)) {
        return;
      }
      if (form.classList.contains("is-loading")) return;
      try {
        form.classList.add("is-loading");
        const formData = new FormData();
        formData.append(textarea.name, textarea.value);

        // Add any other form fields if present
        form.querySelectorAll("input").forEach((input) => {
          if (input.name) {
            formData.append(input.name, input.value);
          }
        });
        const response = await POST(
          `${this.repoLink}/discussions/${this.discussionId}/comment`,
          { data: formData },
        );

        if (!response.ok) {
          return;
        }
        const body = await response.json();

        const resp = await GET(
          `${this.repoLink}/discussions/comment/${body.id}`,
        );
        const commentHolderHTML = await resp.text();

        const tempDiv = document.createElement("div");
        tempDiv.innerHTML = commentHolderHTML;
        const commentHolder = tempDiv.firstElementChild;

        form
          .closest(".discussion-file-comment-holder")
          .replaceWith(commentHolder);
      } catch (e) {
        console.log(e.errorText);
      } finally {
        form.classList.remove("is-loading");
      }
    },

    async fetchDiscussionComment() {
      const codeBlocks = this.content.codeBlocks;

      try {
        const allHtmlTexts = await Promise.all(
          codeBlocks.map(async (codeBlock) => {
            const comments = codeBlock.comments;
            const codeId = codeBlock.codeId;

            const htmlTexts = await Promise.all(
              comments.map(async (comment) => {
                const response = await GET(
                  `${this.repoLink}/discussions/comment/${comment.id}`,
                );
                const result = await response.text();

                const tempDiv = document.createElement("div");
                tempDiv.innerHTML = result;
                const commentHolder = tempDiv.firstElementChild;

                return { comment, commentHolder, codeId };
              }),
            );

            return htmlTexts;
          }),
        );

        allHtmlTexts.flat().forEach(({ comment, commentHolder, codeId }) => {
          const targetLine = this.$refs.codeTable.querySelector(
            `#line-${codeId}-${comment.endLine}`,
          );

          if (targetLine) {
            const tr = document.createElement("tr")
            const td = document.createElement("td")
            td.setAttribute("colspan", "3")
            td.appendChild(commentHolder) 
            tr.appendChild(td)
            targetLine.insertAdjacentElement('afterend', tr);
          }
        });
      } catch (e) {
        console.error("Error processing code blocks:", e);
      }

    },

    async initDiscussionFileCommentForm(form) {
      form.addEventListener("submit", this.submitDiscussionFileCommentForm);
    },

  },
};
</script>

<template>
  <div
    class="file-header ui top attached header tw-items-center tw-justify-between tw-flex-wrap"
    style="position: sticky; top: 0; z-index: 999"
  >
    <div class="file-info tw-font-mono">
      <div :href="`#discussion-${content.NameHash}`" class="file-info-entry">
        {{ content.Name }}
      </div>
    </div>
  </div>
  <div class="ui bottom attached table unstackable segment">
    <div class="file-view code-view" style="display: flex">
      <table :id="content.Name" ref="codeTable" class="discussion-file-table">
        <tbody v-for="codeBlock in content.codeBlocks">
          <tr
            v-for="line in codeBlock.lines"
            :id="`line-${codeBlock.codeId}-${line.lineNumber}`"
            class="code-line"
            :key="`${codeBlock.codeId}-${line.lineNumber}`"
          >
            <td
              class="lines-num"
              :id="`num-${codeBlock.codeId}-${line.lineNumber}`"
            >
              {{ line.lineNumber }}
            </td>
            <td>
              <button
                @mousedown="handleMouseDown"
                @click="showCommentForm"
                class="ui primary button add-code-comment add-code-comment-right"
              >
                <SvgIcon name="octicon-plus" />
              </button>
            </td>
            <td class="lines-code chroma">
              {{ line.content }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
/* 선택된 라인의 스타일을 정의합니다. */
.selected-line {
  background-color: #f5f5dc;
}
/* 테이블이 선택 중일 때의 스타일을 정의합니다. */
.is-selecting {
  cursor: pointer;
}
</style>
