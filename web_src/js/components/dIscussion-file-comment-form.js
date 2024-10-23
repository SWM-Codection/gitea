import { initComboMarkdownEditor, validateTextareaNonEmpty } from "../features/comp/ComboMarkdownEditor.js";
import { initDiscussionCommentEventHandler, initDiscussionCommentsEventHandler } from "../features/discussion-file-comment.js";
import { GET, POST } from "../modules/fetch.js";
import { hideElem, showElem } from "../utils/dom.js";


const { pageData } = window.config;


export async function createCommentPlaceHolder(commentText) {
  const placeholder = document.createElement("tr");
  const td = document.createElement("td");
  td.innerHTML = commentText;
  td.setAttribute("colspan", "3");
  placeholder.appendChild(td);
  await initComboMarkdownEditor(td.querySelector(".combo-markdown-editor"));

  return placeholder;
}

export async function initDiscussionFileCommentForm(form) {
  form.addEventListener("submit", submitDiscussionFileCommentForm);
}

async function submitDiscussionFileCommentForm(event) {
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
    const postUrl = form.getAttribute("action")
    const getUrl = form.getAttribute("data-get-url")
    const formData = new FormData(form);


    const response = await POST(
      postUrl,
      { data: formData },
    );

    if (!response.ok) {
      return;
    }

    const body = await response.json();

    const resp = await GET(
      `${getUrl}/${body.id}`,
    );
    const commentHolderText = await resp.text();

    const commentHolder = convertTextToHTML(commentHolderText)

    form
      .replaceWith(commentHolder);
    initDiscussionCommentsEventHandler(commentHolder)

  } catch (e) {

    console.error(e.message);
  } finally {
    form.classList.remove("is-loading");
  }
}

export function removeCommentForm(event) {
  if (
    event.target &&
    event.target.classList.contains("cancel-code-comment")
  ) {
    const commentForm = event.target.closest("tr");
    if (commentForm) {
      commentForm.remove();
    }
  }
}

function removeEditCommentForm(event) {
  if (
    event.target &&
    event.target.classList.contains("cancel-code-comment")
  ) {
    const commentForm = event.target.closest("tr");
    if (commentForm) {
      commentForm.remove();
    }
  }
}

export function convertTextToHTML(text) {
  const tempDiv = document.createElement("div");
  tempDiv.innerHTML = text;
  return tempDiv.firstElementChild;
}

export async function fetchCommentForm(requestURL) {

  try {
    const response = await GET(requestURL);
    if (!response.ok) {
      
      return;
    }
    const body = await response.text();

    const placeholder = await createCommentPlaceHolder(body);

    placeholder.addEventListener("click", removeCommentForm, {
      capture: true,
    });

    return placeholder

  } catch (err) {
    console.error(err.message);
  }
}

export async function fetchCommentReplyForm(requestURL) {

  try {
    const response = await GET(requestURL);
    if (!response.ok) {
      return;
    }
    const body = await response.text();
    const placeholder = convertTextToHTML(body);
    await initComboMarkdownEditor(placeholder.querySelector(".combo-markdown-editor"));

    placeholder.addEventListener("click", removeCommentForm, {
      capture: true,
    });

    return placeholder

  } catch (err) {
    console.error(err.message);
  }
}


export async function renderReplyCommentForm(event) {

  try {
    const button = event.target;
    const commentHolder = button.closest(".discussion-file-comment-holder");
    const commentsList = commentHolder.querySelector(".comments")
    const discussionId = commentHolder.querySelector("input[name='dId']").value;
    const codeId = commentHolder.querySelector("input[name='codeId']").value;
    const startLine = commentHolder.querySelector("input[name='startLine']").value;
    const endLine = commentHolder.querySelector("input[name='endLine']").value;

    const queryParams = {
      discussionId: discussionId,
      codeId : codeId,
      startLine: startLine,
      endLine: endLine,
    };

    const queries = Object.entries(queryParams).map(([key, value]) => `${key}=${value}`).join('&');
    const requestURL = `${pageData.RepoLink}/discussions/comment?${queries}` 

    const commentForm = await fetchCommentReplyForm(requestURL);
    initDiscussionFileCommentForm(commentForm);
    const groupId = commentHolder.querySelector("input[name='groupId']").value;
    const groupIdInput = document.createElement("input");
    groupIdInput.setAttribute("type", "hidden");
    groupIdInput.setAttribute("name", "groupId");
    groupIdInput.setAttribute("value", groupId);
    commentForm.appendChild(groupIdInput);
    commentsList.appendChild(commentForm);


  } catch (err) {
    console.error(err.message);
  }


}



