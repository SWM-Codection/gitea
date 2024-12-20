import { comment } from "postcss";
import { fetchCommentForm, renderReplyCommentForm } from "../components/dIscussion-file-comment-form.js";
import { DELETE, PUT } from "../modules/fetch.js";
import { hideElem, showElem } from "../utils/dom.js";
import {
  getComboMarkdownEditor,
  initComboMarkdownEditor,
} from "./comp/ComboMarkdownEditor.js";
import { toAbsoluteUrl } from "../utils.js";
import { showTemporaryTooltip } from "../modules/tippy.js";
import { clippie } from "clippie";

export function initDiscussionCommentsEventHandler(commentsHolder) {

  const comments = commentsHolder.querySelectorAll(".comment")
  const isReply = () => {
    return comments.length == 0
  }

  if (isReply()) {
    initDiscussionCommentEventHandler(commentsHolder)
    return
  }

  comments.forEach(comment => {
    initDiscussionCommentEventHandler(comment)
  });
  initDiscussionFileCommentSelectHighlight(commentsHolder)
  initDiscussionCommentReply(commentsHolder);


}


function initDiscussionFileCommentSelectHighlight(commentHolder) {

  commentHolder.addEventListener("click", (event) => {
    event.stopPropagation();

    removeSelectedLines();

    const startLine = parseInt(event.currentTarget.querySelector("input[name='startLine']").value, 10);
    const endLine = parseInt(event.currentTarget.querySelector("input[name='endLine']").value, 10);
    const codeId = event.currentTarget.querySelector("input[name='codeId']").value;

    for (let lineNumber = startLine; lineNumber <= endLine; lineNumber++) {
      const lineElement = document.querySelector(`#line-${codeId}-${lineNumber}`);

      if (lineElement) {
        lineElement.classList.add("selected-line");
      }
    }
  });

  document.addEventListener("click", (event) => {
    if (!commentHolder.contains(event.target)) {
      removeSelectedLines();
    }
  });

  function removeSelectedLines() {
    const selectedLines = document.querySelectorAll(".selected-line");
    selectedLines.forEach(line => line.classList.remove("selected-line"));
  }

}

export async function initDiscussionCommentEventHandler(comment) {
  initDiscussionCommentDropDown(comment);
  initReactionDropdown(comment);
  initQuoteComment(comment);
  initDiscussionCommentDelete(comment);
  initDiscussionCommentUpdate(comment);
}

async function initQuoteComment(comment) {
  const {copy_success, copy_error} = window.config.i18n;
  const QuoteButton = comment.querySelector(".discussion-quote-comment");

  if (!QuoteButton) return;

  QuoteButton.addEventListener("click", async (e) => {

    const target = e.target.closest('[data-clipboard-text], [data-clipboard-target]');
    if (!target) return;
  
    e.preventDefault();
  
    let text = target.getAttribute('data-clipboard-text');
    if (!text) {
      text = document.querySelector(target.getAttribute('data-clipboard-target'))?.value;
    }
  
    if (text && target.getAttribute('data-clipboard-text-type') === 'url') {
      text = toAbsoluteUrl(text);
    }
  
    if (text) {
      const success = await clippie(text);
      showTemporaryTooltip(target, success ? copy_success : copy_error);
    }
  });

}

function initDiscussionCommentDelete(comment) {
  const deleteButton = comment.querySelector(".discussion-delete-comment");

  if (!deleteButton) return;

  deleteButton.addEventListener("click", async (event) => {
    event.preventDefault();

    const confirmationMessage = deleteButton.getAttribute("data-locale");
    const deleteUrl = deleteButton.getAttribute("data-url");
    const commentId = deleteButton.getAttribute("data-comment-id");

    if (window.confirm(confirmationMessage)) {
      try {
        const response = await DELETE(deleteUrl);
        if (!response.ok) throw new Error("Failed to delete comment");

        if (comment.parentElement.children.length === 1) {
          comment.closest(".discussion-file-comment-holder").remove();
          return
        }
        comment.remove()

      } catch (error) {
        // TODO: 다국어 에러 메세지 지원
        alert("코멘트 삭제에 실패했습니다.");
        console.error(error);
      }
    }
  });
}

function initDiscussionCommentUpdate(comment) {
  const updateButton = comment.querySelector(".discussion-edit-content");
  const cancelButton = comment.querySelector(".cancel-edit-code-comment");
  const updateSubmitButton = comment.querySelector(".btn-edit-comment");

  const renderContent = comment.querySelector(".render-content");
  const updateForm = comment.querySelector(".edit-comment-form");
  const rawContent = comment.querySelector(".raw-content");

  let comboMarkdownEditor = initComboMarkdownEditor(
    updateForm.querySelector(".combo-markdown-editor"),
  );

  updateButton?.addEventListener("click", (event) => {
    event.preventDefault();

    showElem(updateForm);
    hideElem(renderContent);

    let comboMarkdownEditor = getComboMarkdownEditor(
      updateForm.querySelector(".combo-markdown-editor"),
    );

    if (!comboMarkdownEditor.value()) {
      comboMarkdownEditor.value(rawContent.textContent);
    }
    comboMarkdownEditor.focus();
  });

  cancelButton?.addEventListener("click", (event) => {
    event.preventDefault();

    showElem(renderContent);
    hideElem(updateForm);
  });

  updateForm?.addEventListener("submit", async (event) => {
    event.preventDefault();
    const updateForm = event.target;

    const formData = new FormData(updateForm);

    let comboMarkdownEditor = getComboMarkdownEditor(
      updateForm.querySelector(".combo-markdown-editor"),
    );
    try {
      const response = await PUT(updateForm.getAttribute("action"), {
        data: formData,
      });
      if (!response.ok) {
        throw Error()
      }
      const data = await response.json();

      if (!data.content) {
        renderContent.innerHTML = document.getElementById("no-content").innerHTML;
        rawContent.textContent = "";
      } else {
        renderContent.innerHTML = data.content;
        rawContent.textContent = comboMarkdownEditor.value();
      }
    }
    catch (e) {
        // TODO: 다국어 에러 메세지 지원
      alert("코멘트 수정에 실패했습니다.")
      console.error(e)
    }


    showElem(renderContent);
    hideElem(updateForm);
  });
}

function initReactionDropdown(comment) {
  var reactionDropdown = comment.querySelector(".select-reaction");
  var reactionDropdownMenu = reactionDropdown.querySelector(".menu");
  if (!reactionDropdown || !reactionDropdownMenu) return;

  reactionDropdown.addEventListener("click", function (event) {
    reactionDropdownMenu.classList.toggle("transition");
    reactionDropdownMenu.classList.toggle("visible");
  });

  document.addEventListener("click", function (event) {
    var isClickInside =
    reactionDropdown.contains(event.target) || reactionDropdown.contains(event.target);

    if (!isClickInside) {
      reactionDropdownMenu.classList.remove("transition");
      reactionDropdownMenu.classList.remove("visible");
    }
  });
}

function initDiscussionCommentDropDown(comment) {
  var dropdown = comment.querySelector(".context-dropdown");
  var dropdownMenu = dropdown.querySelector(".menu");

  if (!dropdown || !dropdownMenu) return;

  dropdown.addEventListener("click", function (event) {
    dropdownMenu.classList.toggle("transition");
    dropdownMenu.classList.toggle("visible");
  });

  document.addEventListener("click", function (event) {
    var isClickInside =
      dropdown.contains(event.target) || dropdown.contains(event.target);

    if (!isClickInside) {
      dropdownMenu.classList.remove("transition");
      dropdownMenu.classList.remove("visible");
    }
  });
}

function initDiscussionCommentReply(commentHolder) {
  const replyButton = commentHolder.querySelector(".discussion-file-comment-form-reply")
  if (!replyButton) return;
  replyButton.addEventListener("click", renderReplyCommentForm)

}


