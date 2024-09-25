import { fetchCommentForm } from "../components/dIscussion-file-comment-form.js";
import { DELETE, PUT } from "../modules/fetch.js";
import { hideElem, showElem } from "../utils/dom.js";
import {
  getComboMarkdownEditor,
  initComboMarkdownEditor,
} from "./comp/ComboMarkdownEditor.js";

export function initDiscussionCommentEventHandler(comment) {
  initDiscussionCommentDropDown(comment);
  initDiscussionCommentDelete(comment);
  initDiscussionCommentUpdate(comment);
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

        const commentElement = comment.querySelector(`#${commentId}`);
        commentElement?.parentElement.remove();
      } catch (error) {
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

  updateButton.addEventListener("click", (event) => {
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

  cancelButton.addEventListener("click", (event) => {
    event.preventDefault();

    showElem(renderContent);
    hideElem(updateForm);
  });

  updateForm.addEventListener("submit", async (event) => {
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
        alert("수정에 실패했습니다.")
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
      console.log(e)
    } 


    showElem(renderContent);
    hideElem(updateForm);
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
