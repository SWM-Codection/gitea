import { DELETE } from "../modules/fetch.js";

export function initDiscussionCommentDelete() {
  const deleteButtons = document.querySelectorAll('.discussion-delete-comment');
  console.log('searching delete comments...')
  if (!deleteButtons) return;
  console.log('find delete buttons done!');
  deleteButtons.forEach((deleteButton) => {
    deleteButton.addEventListener('click', async (event) => {
      event.preventDefault();
  
      const confirmationMessage = deleteButton.getAttribute('data-locale');
      const deleteUrl = deleteButton.getAttribute('data-url');
      const commentId = deleteButton.getAttribute('data-comment-id');
  
      if (window.confirm(confirmationMessage)) {
        try {
          const response = await DELETE(deleteUrl);
          if (!response.ok) throw new Error('Failed to delete comment');
  
          const commentElement = document.querySelector(`#${commentId}`);
          commentElement?.parentElement.remove();
        } catch (error) {
          console.error(error);
        }
      }
    });
  })
}
