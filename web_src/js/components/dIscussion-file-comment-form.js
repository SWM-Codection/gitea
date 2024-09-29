import { initComboMarkdownEditor, validateTextareaNonEmpty } from "../features/comp/ComboMarkdownEditor.js";
import { initDiscussionCommentEventHandler } from "../features/discussion-file-comment.js";
import { GET, POST } from "../modules/fetch.js";

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
    initDiscussionCommentEventHandler(commentHolder)

    form
      .closest(".discussion-file-comment-holder")
      .replaceWith(commentHolder);
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
    const response = await GET(requestURL.toString());
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

export async function fetchEditCommentForm() {

  try {
    const response = await GET()
  } catch(err) {

  }

}



