import { fetchCommentForm } from "../components/dIscussion-file-comment-form.js";
import { DELETE } from "../modules/fetch.js";

export function initDiscussionCommentEventHandler(comment) {
  initDiscussionCommentDropDown(comment);
  initDiscussionCommentDelete(comment);
}

function initDiscussionCommentDelete(comment) {
  const deleteButton = comment.querySelector('.discussion-delete-comment');

  if (!deleteButton) return;

  deleteButton.addEventListener('click', async (event) => {
    event.preventDefault();

    const confirmationMessage = deleteButton.getAttribute('data-locale');
    const deleteUrl = deleteButton.getAttribute('data-url');
    const commentId = deleteButton.getAttribute('data-comment-id');

    if (window.confirm(confirmationMessage)) {
      try {
        const response = await DELETE(deleteUrl);
        if (!response.ok) throw new Error('Failed to delete comment');

        const commentElement = comment.querySelector(`#${commentId}`);
        commentElement?.parentElement.remove();
      } catch (error) {
        console.error(error);
      }
    }
  });
}


function initDiscussionCommentUpdate(comment) {
  const updateButton = comment.querySelector('.discussion-update-comment');

  if (!updateButton) return;

  updateButton.addEventListener('click', async (event) => {
    event.preventDefault();
    const updateUrl = updateButton.getAttribute('data-url');

    fetchCommentForm()

  });
}


function initDiscussionCommentDropDown(comment) {
    
    // 드롭다운 버튼과 메뉴 요소를 가져옵니다.

    var dropdown = comment.querySelector('.context-dropdown');
    var dropdownMenu = dropdown.querySelector(".menu")

    if (!dropdown || !dropdownMenu) return;

    // 드롭다운 버튼 클릭 이벤트 리스너 추가
    dropdown.addEventListener('click', function(event) {

        // 메뉴의 클래스 리스트를 토글하여 드롭다운 열기/닫기 구현
        dropdownMenu.classList.toggle('transition');
        dropdownMenu.classList.toggle('visible');
    });

    // 드롭다운 외부를 클릭하면 드롭다운 닫기
    document.addEventListener('click', function(event) {
        var isClickInside = dropdown.contains(event.target) || dropdown.contains(event.target);

        if (!isClickInside) {
            dropdownMenu.classList.remove('transition');
            dropdownMenu.classList.remove('visible');
        }
    });

}