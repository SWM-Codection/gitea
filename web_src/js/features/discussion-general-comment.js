import { DELETE, POST } from "../modules/fetch.js";

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

export function initDiscussionCommentReaction() { 
  const $reactionButtons = document.querySelectorAll('.discussion-comment-reaction-button');
  $reactionButtons.forEach(($el) => {
    $el.onclick = async (e) => {
      const reactionType = $el.getAttribute('data-tooltip-content');
      const dataUrl = $el.getAttribute('data-url');
      const reacted = $el.getAttribute('data-has-reacted');
      const reactUrl = `${dataUrl}/${reacted ? 'unreact' : 'react'}`;

      const resp = await POST(reactUrl, {data: {'Content': reactionType}}); 

      const data = await resp.json();
      const html = data.html; 
      if (!html) return; 

      // handle reaction add 
      const $commentContainer = $el.closest('.comment-container');
      const $bottomReactions = $commentContainer.querySelector('.bottom-reactions');
      $bottomReactions?.remove(); 
      $commentContainer.insertAdjacentHTML('beforeend', html);
      initDiscussionCommentReaction();
    };
  });
}
